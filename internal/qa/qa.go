package qa

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func LintCommands(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	errors = append(errors, testShortEndWithDotError(commands)...)
	errors = append(errors, testShortIsNotPresentError(commands)...)
	errors = append(errors, testWellKnownArgAtTheEndError(commands)...)
	errors = append(errors, testArgMustUseDashError(commands)...)
	errors = append(errors, testPositionalArgMustBeRequiredError(commands)...)
	errors = append(errors, testExampleCanHaveOnlyOneTypeOfExampleError(commands)...)
	errors = append(errors, testDifferentLocalizationForNamespaceError(commands)...)
	errors = append(errors, testDuplicatedCommandError(commands)...)
	return errors
}

type ShortMustNotEndWithDotError struct {
	Command *core.Command
}

func (err ShortMustNotEndWithDotError) Error() string {
	return "short must not end with '.' for command '" + err.Command.GetCommandLine() + "'"
}

func testShortEndWithDotError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		if strings.HasSuffix(command.Short, ".") {
			errors = append(errors, &ShortMustNotEndWithDotError{
				Command: command,
			})
		}
	}
	return errors
}

type ShortMustBePresentError struct {
	Command *core.Command
}

func (err ShortMustBePresentError) Error() string {
	return "short must be present for command '" + err.Command.GetCommandLine() + "'"
}

func testShortIsNotPresentError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		if command.Short == "" {
			errors = append(errors, &ShortMustBePresentError{
				Command: command,
			})
		}
	}
	return errors
}

type WellKnownArgOrderError struct {
	Command *core.Command
	Argspec *core.ArgSpec
}

func (err WellKnownArgOrderError) Error() string {
	return "well-known arg order must be respected '" + err.Command.GetCommandLine() + "', arg '" + err.Argspec.Name + "'"
}

type WellKnownArgAtTheEndError struct {
	Command *core.Command
	Argspec *core.ArgSpec
}

func (err WellKnownArgAtTheEndError) Error() string {
	return "well-known arg must be at the end'" + err.Command.GetCommandLine() + "', arg '" + err.Argspec.Name + "'"
}

const (
	wellKnownArgOrganizationID = "organization-id"
	wellKnownArgOrganization   = "organization"
	wellKnownArgRegion         = "region"
	wellKnownArgZone           = "zone"
)

var (
	wellKnownArgs = map[string]struct{}{
		wellKnownArgOrganizationID: {},
		wellKnownArgOrganization:   {},
		wellKnownArgRegion:         {},
		wellKnownArgZone:           {},
	}
	wellKnownArgsOrder = []string{
		wellKnownArgOrganizationID,
		wellKnownArgOrganization,
		wellKnownArgRegion,
		wellKnownArgZone,
	}
)

func testWellKnownArgAtTheEndError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		wkaCounter := 0
		wkaNotAtTheEnd := false
		lastWKA := (*core.ArgSpec)(nil)
		for argPosition, argspec := range command.ArgSpecs {
			if _, ok := wellKnownArgs[argspec.Name]; ok {
				respectOrder := false
				for ; wkaCounter < len(wellKnownArgsOrder); wkaCounter++ {
					if argspec.Name == wellKnownArgsOrder[wkaCounter] {
						respectOrder = true
						wkaCounter++ // next well-known arg can't be the same
						break
					}
				}

				if !respectOrder {
					errors = append(errors, &WellKnownArgOrderError{
						Command: command,
						Argspec: argspec,
					})
				}

				wkaNotAtTheEnd = argPosition+1 != len(command.ArgSpecs)
				lastWKA = argspec
			}
		}
		if wkaNotAtTheEnd {
			errors = append(errors, &WellKnownArgAtTheEndError{
				Command: command,
				Argspec: lastWKA,
			})
		}
	}
	return errors
}

type ArgMustUseDashError struct {
	Command *core.Command
	Argspec *core.ArgSpec
}

func (err ArgMustUseDashError) Error() string {
	return "arg must use dash for command '" + err.Command.GetCommandLine() + "', arg '" + err.Argspec.Name + "'"
}

func testArgMustUseDashError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		for _, argspec := range command.ArgSpecs {
			if strings.Contains(argspec.Name, "_") {
				errors = append(errors, &ArgMustUseDashError{
					Command: command,
					Argspec: argspec,
				})
			}
		}
	}
	return errors
}

type PositionalArgMustBeRequiredError struct {
	Command *core.Command
	Argspec *core.ArgSpec
}

func (err PositionalArgMustBeRequiredError) Error() string {
	return "positional argument must be required '" + err.Command.GetCommandLine() + "', arg '" + err.Argspec.Name + "'"
}

func testPositionalArgMustBeRequiredError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		for _, argSpec := range command.ArgSpecs {
			if argSpec.Positional && !argSpec.Required {
				errors = append(errors, &PositionalArgMustBeRequiredError{
					Command: command,
					Argspec: argSpec,
				})
			}
		}
	}
	return errors
}

type ExampleCanHaveOnlyOneTypeOfExampleError struct {
	Command      *core.Command
	ExampleIndex int
}

func (err ExampleCanHaveOnlyOneTypeOfExampleError) Error() string {
	return "arg must use dash for command '" + err.Command.GetCommandLine() + "', example #" + strconv.Itoa(err.ExampleIndex)
}

func testExampleCanHaveOnlyOneTypeOfExampleError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for _, command := range commands.GetAll() {
		for i, example := range command.Examples {
			if example.Request != "" && example.Raw != "" {
				errors = append(errors, &ExampleCanHaveOnlyOneTypeOfExampleError{
					Command:      command,
					ExampleIndex: i,
				})
			}
		}
	}
	return errors
}

type DifferentLocalizationForNamespaceError struct {
	Command1  *core.Command
	Command2  *core.Command
	Checks    []bool
	ArgNames1 []string
	ArgNames2 []string
}

func (err DifferentLocalizationForNamespaceError) Error() string {
	return fmt.Sprintf("different localization for commands '%v', '%v': %v, %v",
		err.Command1.GetCommandLine(), err.Command2.GetCommandLine(), err.ArgNames1, err.ArgNames2)
}

func testDifferentLocalizationForNamespaceError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	for i, command1 := range commands.GetAll() {
		for j, command2 := range commands.GetAll() {
			if i >= j {
				continue
			}

			samePathLength := strings.Count(command1.GetCommandLine(), " ") == strings.Count(command2.GetCommandLine(), " ")

			sameNamespace := command1.Namespace == command2.Namespace

			c1HasRegionOnly := command1.ArgSpecs.GetByName("region") != nil && command1.ArgSpecs.GetByName("zone") == nil
			c2HasRegionOnly := command2.ArgSpecs.GetByName("region") != nil && command2.ArgSpecs.GetByName("zone") == nil

			c1HasZoneOnly := command1.ArgSpecs.GetByName("region") == nil && command1.ArgSpecs.GetByName("zone") != nil
			c2HasZoneOnly := command2.ArgSpecs.GetByName("region") == nil && command2.ArgSpecs.GetByName("zone") != nil

			c1NoRegionNoZone := command1.ArgSpecs.GetByName("region") == nil && command1.ArgSpecs.GetByName("zone") == nil
			c2NoRegionNoZone := command2.ArgSpecs.GetByName("region") == nil && command2.ArgSpecs.GetByName("zone") == nil

			if !samePathLength {
				continue
			}

			if !sameNamespace {
				continue
			}

			switch {
			case c1HasRegionOnly && c2HasRegionOnly:
				// do nothing
			case c1HasZoneOnly && c2HasZoneOnly:
				// do nothing
			case c1NoRegionNoZone && c2NoRegionNoZone:
				// do nothing
			default:

				argNames1 := []string(nil)
				for _, argSpec := range command1.ArgSpecs {
					if argSpec.Name == "zone" || argSpec.Name == "region" {
						argNames1 = append(argNames1, argSpec.Name)
					}
				}

				argNames2 := []string(nil)
				for _, argSpec := range command2.ArgSpecs {
					if argSpec.Name == "zone" || argSpec.Name == "region" {
						argNames2 = append(argNames2, argSpec.Name)
					}
				}

				errors = append(errors, &DifferentLocalizationForNamespaceError{
					Command1:  command1,
					Command2:  command2,
					ArgNames1: argNames1,
					ArgNames2: argNames2,
					Checks:    []bool{sameNamespace, c1HasRegionOnly, c2HasRegionOnly, c1HasZoneOnly, c2HasZoneOnly, c1NoRegionNoZone, c2NoRegionNoZone},
				})
			}
		}
	}
	return errors
}

type DuplicatedCommandError struct {
	Command *core.Command
}

func (err DuplicatedCommandError) Error() string {
	return fmt.Sprintf("duplicated command '%s'", err.Command.GetCommandLine())
}

// testDuplicatedCommandError testes that there is no duplicate command.
func testDuplicatedCommandError(commands *core.Commands) []interface{} {
	errors := []interface{}(nil)
	uniqueness := make(map[string]bool)

	for _, command := range commands.GetAll() {
		key := command.GetCommandLine()

		if uniqueness[key] {
			errors = append(errors, &DuplicatedCommandError{Command: command})
			continue
		}

		uniqueness[key] = true
	}

	return errors
}