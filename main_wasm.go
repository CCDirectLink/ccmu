package main

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/CCDirectLink/CCUpdaterCLI/game"
	"github.com/CCDirectLink/CCUpdaterCLI/pkg"
)

func main() {
	go initialize()
	select {}
}

func initialize() {
	js.Global().Set("game", map[string]interface{}{
		//"Default": mapGame(&game.Default),
		"At": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return mapGame(game.At(args[0].String()))
		}),
		"Default": mapGame(game.Default),
	})

	js.Global().Set("pkg", map[string]interface{}{
		"ErrNotFound": mapError(pkg.ErrNotFound),

		"ModeUnknown":   int(pkg.ModeUnknown),
		"ModeInstall":   int(pkg.ModeInstall),
		"ModeUninstall": int(pkg.ModeUninstall),
		"ModeUpdate":    int(pkg.ModeUpdate),

		"ReasonUnknown":          int(pkg.ReasonUnknown),
		"ReasonNoInternet":       int(pkg.ReasonNoInternet),
		"ReasonAlreadyInstalled": int(pkg.ReasonAlreadyInstalled),
		"ReasonInvalidFormat":    int(pkg.ReasonInvalidFormat),
		"ReasonNotFound":         int(pkg.ReasonNotFound),
		"ReasonNotAvailable":     int(pkg.ReasonNotAvailable),
		"ReasonAccess":           int(pkg.ReasonAccess),
		"ReasonDependant":        int(pkg.ReasonDependant),
	})
}

func mapGame(g game.Game) js.Value {
	return js.ValueOf(map[string]interface{}{
		"Available": promise(func(this js.Value, args []js.Value) interface{} {
			pkgs, err := g.Available()
			return []interface{}{mapPackages(pkgs), mapError(err)}
		}),
		"Installed": promise(func(this js.Value, args []js.Value) interface{} {
			pkgs, err := g.Installed()
			return []interface{}{mapPackages(pkgs), mapError(err)}
		}),

		"Find": promise(func(this js.Value, args []js.Value) interface{} {
			return mapPackages(g.Find(args[0].String()))
		}),
		"Get": promise(func(this js.Value, args []js.Value) interface{} {
			p, err := g.Get(args[0].String())
			return []interface{}{mapPackage(p), mapError(err)}
		}),

		"BasePath": promise(func(this js.Value, args []js.Value) interface{} {
			path, err := g.BasePath()
			return []interface{}{path, mapError(err)}
		}),
	})
}

func mapPackage(p pkg.Package) js.Value {
	if p == nil {
		return js.Undefined()
	}

	return js.ValueOf(map[string]interface{}{
		"Info": promise(func(this js.Value, args []js.Value) interface{} {
			info, err := p.Info()
			return []interface{}{mapInfo(info), mapError(err)}
		}),
		"Installed": promise(func(this js.Value, args []js.Value) interface{} {
			return p.Installed()
		}),
		"Available": promise(func(this js.Value, args []js.Value) interface{} {
			return p.Available()
		}),

		"Install": promise(func(this js.Value, args []js.Value) interface{} {
			return mapError(p.Install())
		}),
		"Uninstall": promise(func(this js.Value, args []js.Value) interface{} {
			return mapError(p.Uninstall())
		}),
		"Update": promise(func(this js.Value, args []js.Value) interface{} {
			return mapError(p.Update())
		}),

		"Dependencies": promise(func(this js.Value, args []js.Value) interface{} {
			pkgs, err := p.Dependencies()
			return []interface{}{mapPackages(pkgs), mapError(err)}
		}),
		"NewestDependencies": promise(func(this js.Value, args []js.Value) interface{} {
			pkgs, err := p.NewestDependencies()
			return []interface{}{mapPackages(pkgs), mapError(err)}
		}),
	})
}

func mapInfo(info pkg.Info) js.Value {
	return js.ValueOf(map[string]interface{}{
		"Name":           info.Name,
		"NiceName":       info.NiceName,
		"Description":    info.Description,
		"Licence":        info.Licence,
		"CurrentVersion": info.CurrentVersion,
		"NewestVersion":  info.NewestVersion,
		"Hidden":         info.Hidden,
		"Outdated": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			res, err := info.Outdated()
			return []interface{}{res, mapError(err)}
		}),
	})
}

func mapPackages(p []pkg.Package) js.Value {
	res := make([]interface{}, len(p))
	for i, pk := range p {
		res[i] = mapPackage(pk)
	}
	return js.ValueOf(res)
}

func mapError(e error) js.Value {
	if e == nil {
		return js.Undefined()
	}
	var pkgErr pkg.Error
	if errors.As(e, &pkgErr) {
		return mapPkgError(pkgErr)
	}
	return js.ValueOf(map[string]interface{}{
		"Error": promise(func(this js.Value, args []js.Value) interface{} {
			return e.Error()
		}),
	})
}

func mapPkgError(e pkg.Error) js.Value {
	return js.ValueOf(map[string]interface{}{
		"Reason": int(e.Reason),
		"Mode":   int(e.Mode),
		"Pkg":    mapPackage(e.Pkg),
		"Err":    mapError(e.Err),
		"Error": promise(func(this js.Value, args []js.Value) interface{} {
			return e.Error()
		}),
		"String": promise(func(this js.Value, args []js.Value) interface{} {
			return e.String()
		}),
		"Unwrap": promise(func(this js.Value, args []js.Value) interface{} {
			return mapError(e.Unwrap())
		}),
	})
}

func promise(callback func(cbThis js.Value, cbArgs []js.Value) interface{}) js.Func {
	return js.FuncOf(func(origThis js.Value, origArgs []js.Value) interface{} {
		return js.Global().Get("Promise").New(js.FuncOf(func(promThis js.Value, promArgs []js.Value) interface{} {
			resolve := promArgs[0]
			reject := promArgs[1]
			go func() {
				defer func() {
					if r := recover(); r != nil {
						reject.Invoke(fmt.Sprintf("%v", r))
					}
				}()
				resolve.Invoke(callback(origThis, origArgs))
			}()
			return nil
		}))
	})
}
