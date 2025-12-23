package main

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

var (
	ole32                = syscall.NewLazyDLL("ole32.dll")
	procCoCreateInstance = ole32.NewProc("CoCreateInstance")
)

type AudioDevice struct {
	ID          string
	Name        string
	DeviceState uint32
}

func getAudioDevices() ([]AudioDevice, error) {
	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return nil, err
	}
	defer ole.CoUninitialize()

	var mmde *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return nil, err
	}
	defer mmde.Release()

	var mmdec *wca.IMMDeviceCollection
	if err := mmde.EnumAudioEndpoints(wca.ERender, wca.DEVICE_STATEMASK_ALL, &mmdec); err != nil {
		return nil, err
	}
	defer mmdec.Release()

	var count uint32
	if err := mmdec.GetCount(&count); err != nil {
		return nil, err
	}

	devices := make([]AudioDevice, 0, count)
	for i := uint32(0); i < count; i++ {
		device, ok := getDeviceInfo(mmdec, i)
		if ok {
			devices = append(devices, device)
		}
	}

	return devices, nil
}

func getDeviceInfo(collection *wca.IMMDeviceCollection, index uint32) (AudioDevice, bool) {
	var mmd *wca.IMMDevice
	if collection.Item(index, &mmd) != nil {
		return AudioDevice{}, false
	}
	defer mmd.Release()

	var deviceID string
	var state uint32
	if mmd.GetId(&deviceID) != nil || mmd.GetState(&state) != nil {
		return AudioDevice{}, false
	}

	var ps *wca.IPropertyStore
	if mmd.OpenPropertyStore(wca.STGM_READ, &ps) != nil {
		return AudioDevice{}, false
	}
	defer ps.Release()

	var pv wca.PROPVARIANT
	if ps.GetValue(&wca.PKEY_Device_FriendlyName, &pv) != nil {
		return AudioDevice{}, false
	}
	defer pv.Clear()

	return AudioDevice{
		ID:          deviceID,
		Name:        pv.String(),
		DeviceState: state,
	}, true
}

func switchAudio(deviceID string) error {
	if deviceID == "" {
		return nil
	}

	err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	if err != nil {
		return err
	}
	defer ole.CoUninitialize()

	policyConfigCLSID := ole.NewGUID("{870af99c-171d-4f9e-af0d-e63df40c2bc9}")
	policyConfigIID := ole.NewGUID("{f8679f50-850a-41cf-9c72-430f290290c8}")

	var policyConfig *IPolicyConfig
	hr, _, _ := procCoCreateInstance.Call(
		uintptr(unsafe.Pointer(policyConfigCLSID)),
		0,
		1,
		uintptr(unsafe.Pointer(policyConfigIID)),
		uintptr(unsafe.Pointer(&policyConfig)),
	)
	if hr != 0 {
		return ole.NewError(hr)
	}
	defer policyConfig.Release()

	deviceIDUTF16, err := syscall.UTF16PtrFromString(deviceID)
	if err != nil {
		return err
	}

	for role := 0; role < 3; role++ {
		ret := policyConfig.SetDefaultEndpoint(deviceIDUTF16, role)
		if ret != 0 {
			return ole.NewError(ret)
		}
	}

	return nil
}

type IPolicyConfig struct {
	vtbl *IPolicyConfigVtbl
}

type IPolicyConfigVtbl struct {
	QueryInterface        uintptr
	AddRef                uintptr
	Release               uintptr
	GetMixFormat          uintptr
	GetDeviceFormat       uintptr
	ResetDeviceFormat     uintptr
	SetDeviceFormat       uintptr
	GetProcessingPeriod   uintptr
	SetProcessingPeriod   uintptr
	GetShareMode          uintptr
	SetShareMode          uintptr
	GetPropertyValue      uintptr
	SetPropertyValue      uintptr
	SetDefaultEndpoint    uintptr
	SetEndpointVisibility uintptr
}

func (pc *IPolicyConfig) Release() uintptr {
	ret, _, _ := syscall.SyscallN(
		pc.vtbl.Release,
		uintptr(unsafe.Pointer(pc)),
	)
	return ret
}

func (pc *IPolicyConfig) SetDefaultEndpoint(deviceID *uint16, role int) uintptr {
	ret, _, _ := syscall.SyscallN(
		pc.vtbl.SetDefaultEndpoint,
		uintptr(unsafe.Pointer(pc)),
		uintptr(unsafe.Pointer(deviceID)),
		uintptr(role),
	)
	return ret
}
