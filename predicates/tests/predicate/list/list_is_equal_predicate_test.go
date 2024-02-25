// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package tests

import (
	"testing"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/predicates/predicate/list"
	"github.com/conformize/conformize/predicates/predicatefactory"
	"github.com/conformize/conformize/predicates/tests"
)

func TestListIsEqualPredicate(t *testing.T) {
	tests := []struct {
		name    string
		value   typed.Valuable
		args    *typed.TupleValue
		want    bool
		wantErr bool
	}{
		{
			name: "returns true when list contain same elements",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{}),
				}, &typed.StringTyped{},
			),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{}),
						}, &typed.StringTyped{},
					),
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.StringTyped{}}},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when lists don't contain the same elements",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{}),
				}, &typed.StringTyped{}),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{}),
						}, &typed.StringTyped{},
					),
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.StringTyped{}}},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false when lists don't contain the same elements at beginning",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal("c", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{}),
				}, &typed.StringTyped{},
			),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{}),
						}, &typed.StringTyped{},
					),
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.StringTyped{}}},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false when lists don't contain the same elements at end",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal("c", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("a", &typed.StringTyped{}),
				}, &typed.StringTyped{}),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal("c", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("e", &typed.StringTyped{}),
						}, &typed.StringTyped{},
					),
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.StringTyped{}}},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when nested lists contain same elements",
			value: typed.NewListValue([]typed.Valuable{
				typed.NewListValue(
					[]typed.Valuable{tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{})}, &typed.StringTyped{},
				),
				typed.NewListValue(
					[]typed.Valuable{tests.PrimVal("d", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{})}, &typed.StringTyped{},
				),
			}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					typed.NewListValue([]typed.Valuable{
						typed.NewListValue(
							[]typed.Valuable{tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{})}, &typed.StringTyped{},
						),
						typed.NewListValue(
							[]typed.Valuable{tests.PrimVal("d", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{})}, &typed.StringTyped{},
						),
					}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}),
				},
				ElementsTypes: []typed.Typeable{
					&typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when nested lists don't contain same elements",
			value: typed.NewListValue([]typed.Valuable{
				typed.NewListValue(
					[]typed.Valuable{tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{})}, &typed.StringTyped{},
				),
				typed.NewListValue(
					[]typed.Valuable{tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{})}, &typed.StringTyped{},
				),
			}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					&typed.ListValue{
						Elements: []typed.Valuable{
							typed.NewListValue([]typed.Valuable{
								typed.NewListValue(
									[]typed.Valuable{tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{})}, &typed.StringTyped{},
								),
								typed.NewListValue(
									[]typed.Valuable{tests.PrimVal("c", &typed.StringTyped{}), tests.PrimVal("e", &typed.StringTyped{})}, &typed.StringTyped{},
								),
							}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}),
						},
						ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}},
					},
				},
				ElementsTypes: []typed.Typeable{
					&typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false when nested lists don't contain same elements at beginning",
			value: typed.NewListValue([]typed.Valuable{
				typed.NewListValue(
					[]typed.Valuable{
						tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("a", &typed.StringTyped{}),
					}, &typed.StringTyped{},
				),
				typed.NewListValue(
					[]typed.Valuable{
						tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{}),
					}, &typed.StringTyped{},
				),
			}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					&typed.ListValue{
						Elements: []typed.Valuable{
							typed.NewListValue([]typed.Valuable{
								typed.NewListValue(
									[]typed.Valuable{
										tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("a", &typed.StringTyped{}),
									}, &typed.StringTyped{},
								),
								typed.NewListValue(
									[]typed.Valuable{
										tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("c", &typed.StringTyped{}),
									}, &typed.StringTyped{},
								),
							}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}),
						},
						ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}},
					},
				},
				ElementsTypes: []typed.Typeable{
					&typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false when nested lists don't contain same elements at end",
			value: typed.NewListValue([]typed.Valuable{
				typed.NewListValue(
					[]typed.Valuable{
						tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}),
					}, &typed.StringTyped{},
				),
				typed.NewListValue(
					[]typed.Valuable{
						tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("e", &typed.StringTyped{}),
					}, &typed.StringTyped{},
				),
			}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					&typed.ListValue{
						Elements: []typed.Valuable{
							typed.NewListValue([]typed.Valuable{
								typed.NewListValue(
									[]typed.Valuable{
										tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}),
									}, &typed.StringTyped{},
								),
								typed.NewListValue(
									[]typed.Valuable{
										tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{}), tests.PrimVal("d", &typed.StringTyped{}),
									}, &typed.StringTyped{},
								),
							}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}),
						},
					},
				},
				ElementsTypes: []typed.Typeable{
					&typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name:  "returns false and error when value is not a list",
			value: typed.NewMapValue(map[string]typed.Valuable{}, &typed.BooleanTyped{}),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(42, &typed.NumberTyped{}), tests.PrimVal(42, &typed.NumberTyped{})}, &typed.NumberTyped{},
					),
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and error when nested lists have different element types",
			value: typed.NewListValue(
				[]typed.Valuable{tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{})}, &typed.StringTyped{},
			),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{tests.PrimVal(42, &typed.NumberTyped{}), tests.PrimVal(42, &typed.NumberTyped{})}, &typed.NumberTyped{},
					),
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			},
			want:    false,
			wantErr: true,
		},
		{
			name:  "returns false when nested lists have different number of elements",
			value: typed.NewListValue([]typed.Valuable{tests.PrimVal("a", &typed.StringTyped{})}, &typed.StringTyped{}),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{tests.PrimVal("a", &typed.StringTyped{}), tests.PrimVal("b", &typed.StringTyped{})}, &typed.StringTyped{},
					),
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.StringTyped{}}},
			},
			want:    false,
			wantErr: false,
		},
		{
			name:  "returns true when both lists are mepty",
			value: typed.NewListValue([]typed.Valuable{}, &typed.BooleanTyped{}),
			args: &typed.TupleValue{
				Elements:      []typed.Valuable{typed.NewListValue([]typed.Valuable{}, &typed.BooleanTyped{})},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.BooleanTyped{}}},
			},
			want:    true,
			wantErr: false,
		},
		{
			name:  "returns false and error when argument is not a list",
			value: typed.NewListValue([]typed.Valuable{tests.PrimVal("a", &typed.StringTyped{})}, &typed.StringTyped{}),
			args: &typed.TupleValue{
				Elements:      []typed.Valuable{typed.NewMapValue(map[string]typed.Valuable{}, &typed.BooleanTyped{})},
				ElementsTypes: []typed.Typeable{&typed.MapTyped{ElementsType: &typed.BooleanTyped{}}},
			},
			want:    false,
			wantErr: true,
		},
		{
			name:    "returns false and error when argument is not provided",
			value:   typed.NewListValue([]typed.Valuable{tests.PrimVal("a", &typed.StringTyped{})}, &typed.StringTyped{}),
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false when single element lists are not equal",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal("a", &typed.StringTyped{}),
				}, &typed.StringTyped{},
			),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					typed.NewListValue([]typed.Valuable{tests.PrimVal("c", &typed.StringTyped{})}, &typed.StringTyped{}),
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.StringTyped{}}},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when single element lists are equal",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal("b", &typed.StringTyped{}),
				}, &typed.StringTyped{},
			),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					typed.NewListValue([]typed.Valuable{tests.PrimVal("b", &typed.StringTyped{})}, &typed.StringTyped{}),
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.StringTyped{}}},
			},
			want:    true,
			wantErr: false,
		},
	}

	listEqPrd := &list.ListIsEqualPredicate{PredicateBuilder: predicatefactory.Instance()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listEqPrd.Test(tt.value, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListIsEqualPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ListIsEqualPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListIsEqualPredicate_Arguments(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "returns 1 as number of expected arguments",
			want: 1,
		},
	}
	listEqPrd := &list.ListIsEqualPredicate{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listEqPrd.Arguments(); got != tt.want {
				t.Errorf("ListIsEqualPredicate.Arguments() = %v, want %v", got, tt.want)
			}
		})
	}
}
