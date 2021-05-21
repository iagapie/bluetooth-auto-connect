package app

import (
	"fmt"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"os"
	"os/signal"
	"syscall"
)

const (
	orgBluezPath      = "/org/bluez"
	orgBluezInterface = "org.bluez"

	adapterInterface = "org.bluez.Adapter1"
	deviceInterface  = "org.bluez.Device1"

	getAllPropertiesMethod = "org.freedesktop.DBus.Properties.GetAll"
	deviceConnectMethod    = "org.bluez.Device1.Connect"

	propName      = "Name"
	propConnected = "Connected"
	propTrusted   = "Trusted"
	propPowered   = "Powered"

	connectedMsg  = "Device %s is already connected.\n"
	connectingMsg = "Connecting to device %s.\n"
	trustedMsg    = "Device %s is not trusted.\n"
)

func Run() {
	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}

	connectDevicesForAllAdapters(conn)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, os.Interrupt)

	exit := make(chan int)

	go func() {
		for {
			s := <-sig
			switch s {
			case syscall.SIGHUP:
				connectDevicesForAllAdapters(conn)
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				exit <- 0
			default:
				exit <- 2
			}
		}
	}()

	code := <-exit
	os.Exit(code)
}

func dbusGetChildObjectPaths(conn *dbus.Conn, path dbus.ObjectPath) ([]dbus.ObjectPath, error) {
	paths := make([]dbus.ObjectPath, 0)
	obj := conn.Object(orgBluezInterface, path)

	node, err := introspect.Call(obj)
	if err != nil {
		return paths, err
	}

	for _, c := range node.Children {
		if c.XMLName.Local == "node" {
			paths = append(paths, dbus.ObjectPath(fmt.Sprintf("%s/%s", path, c.Name)))
		}
	}

	return paths, nil
}

func connectDevicesForAdapter(conn *dbus.Conn, adapterPath dbus.ObjectPath) {
	devicePaths, err := dbusGetChildObjectPaths(conn, adapterPath)
	if err != nil {
		panic(err)
	}

	for _, devicePath := range devicePaths {
		obj := conn.Object(orgBluezInterface, devicePath)

		result := make(map[string]dbus.Variant)
		err = obj.Call(getAllPropertiesMethod, 0, deviceInterface).Store(&result)
		if err != nil {
			panic(err)
		}

		name, ok := result[propName]
		if !ok {
			continue
		}

		if result[propConnected].Value().(bool) {
			fmt.Printf(connectedMsg, name.String())
			continue
		}

		if result[propTrusted].Value().(bool) == false {
			fmt.Printf(trustedMsg, name.String())
			continue
		}

		obj.Go(deviceConnectMethod, 0, nil)

		fmt.Printf(connectingMsg, name.String())
	}
}

func connectDevicesForAllAdapters(conn *dbus.Conn) {
	adapterPaths, err := dbusGetChildObjectPaths(conn, orgBluezPath)
	if err != nil {
		panic(err)
	}

	for _, adapterPath := range adapterPaths {
		result := make(map[string]dbus.Variant)

		err = conn.Object(orgBluezInterface, adapterPath).Call(getAllPropertiesMethod, 0, adapterInterface).Store(&result)
		if err != nil {
			panic(err)
		}

		if result[propPowered].Value().(bool) {
			connectDevicesForAdapter(conn, adapterPath)
		}
	}
}
