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
					"bruceBanner": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"robForwardlund": UserPermission{
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
					"bruceBanner": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"robForwardlund": UserPermission{
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
					"bruceBanner": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"robForwardlund": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
				},
			},
			"bruceBanner": AccountInfo{
				LiquidBalance:   7,
				IlliquidBalance: 8,
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
					"thePebble": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"robForwardlund": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
				},
			},
			"robForwardlund": AccountInfo{
				LiquidBalance:   9,
				IlliquidBalance: 10,
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
					"thePebble": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"bruceBanner": UserPermission{
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
			"second": TableInfo{
				Fields: []string{
					"first",
					"second",
					"third",
					"forth",
				},
				Cells: [][]CellLocation{
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(4),
							Position:        big.NewInt(4),
							PostionInRecord: big.NewInt(4),
						},
						CellLocation{
							BlockNumber:     big.NewInt(4),
							Position:        big.NewInt(5),
							PostionInRecord: big.NewInt(5),
						},
						CellLocation{
							BlockNumber:     big.NewInt(4),
							Position:        big.NewInt(6),
							PostionInRecord: big.NewInt(6),
						},
					},
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(5),
							Position:        big.NewInt(4),
							PostionInRecord: big.NewInt(4),
						},
						CellLocation{
							BlockNumber:     big.NewInt(5),
							Position:        big.NewInt(5),
							PostionInRecord: big.NewInt(5),
						},
						CellLocation{
							BlockNumber:     big.NewInt(5),
							Position:        big.NewInt(6),
							PostionInRecord: big.NewInt(6),
						},
					},
				},
				Writes: []RecordLocation{
					RecordLocation{
						BlockNumber: *big.NewInt(4),
						Position:    *big.NewInt(4),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(5),
						Position:    *big.NewInt(5),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(6),
						Position:    *big.NewInt(6),
					},
				},
			},
			"third": TableInfo{
				Fields: []string{
					"first",
					"second",
					"third",
					"forth",
				},
				Cells: [][]CellLocation{
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(7),
							Position:        big.NewInt(7),
							PostionInRecord: big.NewInt(7),
						},
						CellLocation{
							BlockNumber:     big.NewInt(7),
							Position:        big.NewInt(8),
							PostionInRecord: big.NewInt(8),
						},
						CellLocation{
							BlockNumber:     big.NewInt(7),
							Position:        big.NewInt(9),
							PostionInRecord: big.NewInt(9),
						},
					},
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(8),
							Position:        big.NewInt(7),
							PostionInRecord: big.NewInt(7),
						},
						CellLocation{
							BlockNumber:     big.NewInt(8),
							Position:        big.NewInt(8),
							PostionInRecord: big.NewInt(8),
						},
						CellLocation{
							BlockNumber:     big.NewInt(8),
							Position:        big.NewInt(9),
							PostionInRecord: big.NewInt(9),
						},
					},
				},
				Writes: []RecordLocation{
					RecordLocation{
						BlockNumber: *big.NewInt(7),
						Position:    *big.NewInt(7),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(8),
						Position:    *big.NewInt(8),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(9),
						Position:    *big.NewInt(9),
					},
				},
			},
			"forth": TableInfo{
				Fields: []string{
					"first",
					"second",
					"third",
					"forth",
				},
				Cells: [][]CellLocation{
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(10),
							Position:        big.NewInt(10),
							PostionInRecord: big.NewInt(10),
						},
						CellLocation{
							BlockNumber:     big.NewInt(10),
							Position:        big.NewInt(11),
							PostionInRecord: big.NewInt(11),
						},
						CellLocation{
							BlockNumber:     big.NewInt(10),
							Position:        big.NewInt(12),
							PostionInRecord: big.NewInt(12),
						},
					},
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(11),
							Position:        big.NewInt(10),
							PostionInRecord: big.NewInt(10),
						},
						CellLocation{
							BlockNumber:     big.NewInt(11),
							Position:        big.NewInt(11),
							PostionInRecord: big.NewInt(11),
						},
						CellLocation{
							BlockNumber:     big.NewInt(11),
							Position:        big.NewInt(12),
							PostionInRecord: big.NewInt(12),
						},
					},
				},
				Writes: []RecordLocation{
					RecordLocation{
						BlockNumber: *big.NewInt(10),
						Position:    *big.NewInt(10),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(11),
						Position:    *big.NewInt(11),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(12),
						Position:    *big.NewInt(12),
					},
				},
			},
		},
	}

	return *dummyData
}
