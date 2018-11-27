package main

import "math/big"

func dummyDebasedMetaData() DebasedMetadata {
	dummyData := &DebasedMetadata{
		Accounts: map[string]AccountInfo{
			"cantSeeMe": AccountInfo{
				LiquidBalance:   1,
				IlliquidBalance: 2,
				Permissions: map[string]UserPermission{
					"theOvertaker": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"thePebble": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
				},
			},
			"theOvertaker": AccountInfo{
				LiquidBalance:   3,
				IlliquidBalance: 4,
				Permissions: map[string]UserPermission{
					"cantSeeMe": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"thePebble": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
				},
			},
			"thePebble": AccountInfo{
				LiquidBalance:   5,
				IlliquidBalance: 6,
				Permissions: map[string]UserPermission{
					"theOvertaker": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"cantSeeMe": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
				},
			},
		},
		Tables: map[string]TableInfo{
			"first": TableInfo{
				Fields: []string{
					"first",
					"second",
					"third",
					"forth",
				},
				Cells: [][]CellLocation{
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(1),
							Position:        big.NewInt(1),
							PostionInRecord: big.NewInt(1),
						},
						CellLocation{
							BlockNumber:     big.NewInt(1),
							Position:        big.NewInt(2),
							PostionInRecord: big.NewInt(2),
						},
						CellLocation{
							BlockNumber:     big.NewInt(1),
							Position:        big.NewInt(3),
							PostionInRecord: big.NewInt(3),
						},
					},
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(2),
							Position:        big.NewInt(1),
							PostionInRecord: big.NewInt(1),
						},
						CellLocation{
							BlockNumber:     big.NewInt(2),
							Position:        big.NewInt(2),
							PostionInRecord: big.NewInt(2),
						},
						CellLocation{
							BlockNumber:     big.NewInt(2),
							Position:        big.NewInt(3),
							PostionInRecord: big.NewInt(3),
						},
					},
				},
				Writes: []RecordLocation{
					RecordLocation{
						BlockNumber: *big.NewInt(1),
						Position:    *big.NewInt(1),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(2),
						Position:    *big.NewInt(2),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(3),
						Position:    *big.NewInt(3),
					},
				},
			},
		},
	}

	return *dummyData
}
