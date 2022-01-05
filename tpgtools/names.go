package main

import "strings"

// We struggle with many types of names, and many conversions between types of names.  In order to sanitize all of this,
// we have four interfaces, and many string-aliased types, which conform to at most one of the interfaces.
// This way, we can't wind up with complex, lossy, and hard-to-trace name conversions, which used to be a big problem in tpgtools.
// The functions are not really meant to be called, so much, but they are useful for converting back to regular strings where needed.
type snakeCaseName interface {
	snakecase() string
}

type titleCaseName interface {
	titlecase() string
}

type jsonCaseName interface {
	jsoncase() string
}

type lowercaseName interface {
	lowercase() string
}

// e.g. `google_compute_instance` or `google_orgpolicy_policy`.
type SnakeCaseTerraformResourceName string

func (s SnakeCaseTerraformResourceName) snakecase() string {
	return string(s)
}

// e.g. `ComputeInstanceGroupManager`.
type TitleCaseFullName string

func (s TitleCaseFullName) titlecase() string {
	return string(s)
}

// e.g. "compute_firewall_rule"
type SnakeCaseFullName string

func (s SnakeCaseFullName) snakecase() string {
	return string(s)
}

// e.g. "os_policy"
type SnakeCaseProductName string

func (s SnakeCaseProductName) snakecase() string {
	return string(s)
}

func (s SnakeCaseProductName) ToTitle() RenderedString {
	return RenderedString(snakeToTitleCase(s).titlecase())
}

// e.g. "ForwardingRule"
type TitleCaseResourceName string

func (t TitleCaseResourceName) titlecase() string {
	return string(t)
}

// e.g. "computeinstancegroupmanager".
type ConjoinedString string

// snakeToLowercase converts a snake_case string to a conjoined string
func snakeToLowercase(s snakeCaseName) ConjoinedString {
	return ConjoinedString(strings.Join(snakeToParts(s, false), ""))
}

// snakeToTitleCase converts a snake_case string to TitleCase / Go struct case.
func snakeToTitleCase(s snakeCaseName) miscellaneousNameTitleCase {
	return miscellaneousNameTitleCase(strings.Join(snakeToParts(s, true), ""))
}

// A type for a string that is not meant for further conversion.  Some functions return a
// RenderedString to indicate that they have been lossily converted to another format.
type RenderedString string

func (r RenderedString) String() string {
	return string(r)
}

func renderSnakeAsTitle(s snakeCaseName) RenderedString {
	return RenderedString(strings.Join(snakeToParts(s, true), ""))
}

// e.g. "ospolicy"
type DCLPackageName string

func (d DCLPackageName) lowercase() string {
	return string(d)
}

type BasePathOverrideNameSnakeCase string

func (b BasePathOverrideNameSnakeCase) snakecase() string {
	return string(b)
}
func (b BasePathOverrideNameSnakeCase) ToUpper() RenderedString {
	return RenderedString(strings.ToUpper(string(b)))
}

func (b BasePathOverrideNameSnakeCase) ToTitle() RenderedString {
	title := snakeToTitleCase(b).titlecase()
	// Got to special case the capitalization of "OS" in "OSConfig", for base paths specifically,
	// because of interop with MMv1.
	if strings.HasPrefix(string(b), "os") {
		return RenderedString("OS" + title[2:])
	}
	return RenderedString(title)
}

// A path on the filesystem, usually relative to the root of the tpgtools/ directory.
type Filepath string

// A package path, potentially including a version suffix.
// e.g. "ospolicy/beta" or "ospolicy"
type DCLPackageNameWithVersion string

// A type for some string, not one of the things that have a specific type above, which is in
// a particular case.  This is useful because we want to be able to write strings functions that take in
// a snake case string or return a snake case string, which works even if the string isn't a
// specific type.
//
// Also, having all these misc strings prevents us from winding up with a bunch of `string` types
// for things that should be explicit.
type miscellaneousNameSnakeCase string

func (m miscellaneousNameSnakeCase) snakecase() string {
	return string(m)
}

type miscellaneousNameTitleCase string

func (m miscellaneousNameTitleCase) titlecase() string {
	return string(m)
}

type miscellaneousNameLowercase string

func (m miscellaneousNameLowercase) lowercase() string {
	return string(m)
}
