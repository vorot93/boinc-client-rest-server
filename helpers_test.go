package main

import (
	"errors"
	"testing"
)

func TestExpandEmptyTags(t *testing.T) {
	var err error
	compErr := errors.New("Mismatch.")
	s := `<boinc_gui_rpc_request><auth1/></boinc_gui_rpc_request>`
	expectation := `<boinc_gui_rpc_request><auth1></auth1></boinc_gui_rpc_request>`
	result := ExpandEmptyTags(s, []string{"auth1"})

	if len(expectation) != len(result) {
		err = compErr
	} else {
		if result != expectation {
			err = compErr
		}
	}

	if err != nil {
		t.Errorf(ErrorOut(expectation, result))
	}
}

func TestCollapseEmptyTags(t *testing.T) {
	var err error
	compErr := errors.New("Mismatch.")
	s := `<boinc_gui_rpc_request><auth1></auth1></boinc_gui_rpc_request>`
	expectation := `<boinc_gui_rpc_request><auth1/></boinc_gui_rpc_request>`
	result := CollapseEmptyTags(s, []string{"auth1"})

	if len(expectation) != len(result) {
		err = compErr
	} else {
		if result != expectation {
			err = compErr
		}
	}

	if err != nil {
		t.Errorf(ErrorOut(expectation, result))
	}
}
