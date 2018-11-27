package main

import "math/big"

func dummyDebasedMetaData() DebasedMetadata {
	dummyData := &DebasedMetadata{
		Accounts: map[string]AccountInfo{
			"cantSeeMe": AccountInfo{
				LiquidBalance:   1,
				IlliquidBalance: 2,
				Permissions: map[string]UserPermission{
					"first": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"photos": UserPermission{
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
					"first": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"photos": UserPermission{
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
					"first": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"photos": UserPermission{
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
			"photos": TableInfo{
				Fields: []string{
					"jpeg",
					"gif",
					"jif",
				},
				Cells: [][]CellLocation{
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(5),
							Position:        big.NewInt(9),
							PostionInRecord: big.NewInt(23),
						},
						CellLocation{
							BlockNumber:     big.NewInt(5),
							Position:        big.NewInt(8),
							PostionInRecord: big.NewInt(2),
						},
						CellLocation{
							BlockNumber:     big.NewInt(4),
							Position:        big.NewInt(3),
							PostionInRecord: big.NewInt(37),
						},
					},
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(7),
							Position:        big.NewInt(67),
							PostionInRecord: big.NewInt(14),
						},
						CellLocation{
							BlockNumber:     big.NewInt(24),
							Position:        big.NewInt(32),
							PostionInRecord: big.NewInt(3),
						},
						CellLocation{
							BlockNumber:     big.NewInt(4),
							Position:        big.NewInt(36),
							PostionInRecord: big.NewInt(7),
						},
					},
				},
				Writes: []RecordLocation{
					RecordLocation{
						BlockNumber: *big.NewInt(2),
						Position:    *big.NewInt(4),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(4),
						Position:    *big.NewInt(7),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(4),
						Position:    *big.NewInt(9),
					},
				},
			},
		},
	}

	return *dummyData
}
