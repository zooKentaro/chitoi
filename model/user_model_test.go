package model

import (
    "testing"

    "github.com/uenoryo/chitoi/database/row"
)

func TestUser_exhaustStamina(t *testing.T) {
    t.Parallel()

    type Test struct {
        Title         string
        User          *User
        ExpectStamina uint32
        IsError       bool
    }

    tests := []Test{
        {
            Title: "success",
            User: &User{
                Row: &row.User{
                    Stamina: 2,
                },
            },
            ExpectStamina: 1,
            IsError:       false,
        },
        {
            Title: "fail: not enough stamina",
            User: &User{
                Row: &row.User{
                    Stamina: 0,
                },
            },
            IsError: true,
        },
    }

    for _, test := range tests {
        t.Run(test.Title, func(t *testing.T) {
            err := test.User.exhaustStamina()
            if test.IsError {
                if err == nil {
                    t.Error("expected error, but not thrown")
                }
                return
            }
            if err != nil {
                t.Errorf("error exhaust stamina, error: %s", err.Error())
            }

            if g, w := test.User.Row.Stamina, test.ExpectStamina; g != w {
                t.Errorf("error stamina %d, want %d", g, w)
            }
        })
    }
}
