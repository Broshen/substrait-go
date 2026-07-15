// SPDX-License-Identifier: Apache-2.0

package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	substraitgo "github.com/substrait-io/substrait-go/v4"
	"github.com/substrait-io/substrait-go/v4/types"
)

func TestStructFieldRefGetTypeRejectsOutOfBoundsField(t *testing.T) {
	parentType := &types.StructType{Types: []types.Type{
		&types.BooleanType{},
		&types.Int32Type{},
		&types.StringType{},
	}}

	ref := &StructFieldRef{Field: int32(len(parentType.Types))}

	var (
		gotType types.Type
		err     error
	)

	require.NotPanics(t, func() {
		gotType, err = ref.GetType(parentType)
	})
	require.ErrorIs(t, err, substraitgo.ErrInvalidType)
	assert.Nil(t, gotType)
}

func TestStructFieldRefGetTypeRejectsNegativeField(t *testing.T) {
	parentType := &types.StructType{Types: []types.Type{
		&types.BooleanType{},
		&types.Int32Type{},
		&types.StringType{},
	}}

	ref := &StructFieldRef{Field: -1}

	var (
		gotType types.Type
		err     error
	)

	require.NotPanics(t, func() {
		gotType, err = ref.GetType(parentType)
	})
	require.ErrorIs(t, err, substraitgo.ErrInvalidType)
	assert.Nil(t, gotType)
}

func TestStructFieldRefGetTypeReturnsInRangeFieldType(t *testing.T) {
	expected := &types.Int64Type{}
	parentType := &types.StructType{Types: []types.Type{
		&types.BooleanType{},
		&types.Int32Type{},
		expected,
	}}

	ref := &StructFieldRef{Field: int32(len(parentType.Types) - 1)}

	gotType, err := ref.GetType(parentType)
	require.NoError(t, err)
	require.NotNil(t, gotType)
	assert.Truef(t, expected.Equals(gotType), "expected type %s, got %s", expected, gotType)
}

func TestStructFieldRefGetTypeResolvesChildType(t *testing.T) {
	expected := &types.StringType{}
	nested := &types.StructType{Types: []types.Type{
		&types.BooleanType{},
		expected,
	}}
	parentType := &types.StructType{Types: []types.Type{
		&types.Int32Type{},
		nested,
		&types.Int64Type{},
	}}

	ref := &StructFieldRef{Field: 1, Child: &StructFieldRef{Field: 1}}

	gotType, err := ref.GetType(parentType)
	require.NoError(t, err)
	require.NotNil(t, gotType)
	assert.Truef(t, expected.Equals(gotType), "expected type %s, got %s", expected, gotType)
}

func TestStructFieldRefGetTypeRejectsNonStructParent(t *testing.T) {
	ref := &StructFieldRef{Field: 0}

	gotType, err := ref.GetType(&types.ListType{Type: &types.Int32Type{}})
	require.ErrorIs(t, err, substraitgo.ErrInvalidType)
	assert.Nil(t, gotType)
}
