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
					"cantSeeMe": UserPermission{
						Roles: [8]bool{
							false, false, true, true, false, true, false, true,
						},
					},
					"theOverTaker": UserPermission{
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
			"champions": TableInfo{
				Fields: []string{
					"cantSeeMe",
					"theOvertaker",
					"TeddyChu",
				},
				Cells: [][]CellLocation{
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(2),
							Position:        big.NewInt(6),
							PostionInRecord: big.NewInt(1000),
						},
						CellLocation{
							BlockNumber:     big.NewInt(3),
							Position:        big.NewInt(17),
							PostionInRecord: big.NewInt(121),
						},
						CellLocation{
							BlockNumber:     big.NewInt(91),
							Position:        big.NewInt(13),
							PostionInRecord: big.NewInt(29),
						},
					},
					[]CellLocation{
						CellLocation{
							BlockNumber:     big.NewInt(4),
							Position:        big.NewInt(7),
							PostionInRecord: big.NewInt(23),
						},
						CellLocation{
							BlockNumber:     big.NewInt(31),
							Position:        big.NewInt(82),
							PostionInRecord: big.NewInt(43),
						},
						CellLocation{
							BlockNumber:     big.NewInt(12),
							Position:        big.NewInt(6),
							PostionInRecord: big.NewInt(9),
						},
					},
				},
				Writes: []RecordLocation{
					RecordLocation{
						BlockNumber: *big.NewInt(9),
						Position:    *big.NewInt(21),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(12),
						Position:    *big.NewInt(91),
					},
					RecordLocation{
						BlockNumber: *big.NewInt(32),
						Position:    *big.NewInt(56),
					},
				},
			},
		},
	}

	return *dummyData
}
