package main

import (
	"github.com/transcom/milmove_orders/pkg/gen/ordersmessages"
)

var enlistedOrdersTypes = map[string]ordersmessages.OrdersType{
	"1": ordersmessages.OrdersTypeIpcot,          // IPCOT In-place consecutive overseas travel
	"8": ordersmessages.OrdersTypeOteip,          // Overseas Tour Extension Incentive Program (OTEIP)
	"9": ordersmessages.OrdersTypeTraining,       // NAVCAD (Naval Cadet) Training
	"A": ordersmessages.OrdersTypeAccession,      // Accession Travel Recruits
	"B": ordersmessages.OrdersTypeAccession,      // Non-recruit Accession Travel
	"C": ordersmessages.OrdersTypeTraining,       // Training Travel
	"D": ordersmessages.OrdersTypeOperational,    // Operational Travel
	"E": ordersmessages.OrdersTypeSeparation,     // Separation Travel
	"F": ordersmessages.OrdersTypeUnitMove,       // Organized Unit/Homeport Change
	"G": ordersmessages.OrdersTypeAccession,      // Midshipman Accession Travel
	"H": ordersmessages.OrdersTypeSpecialPurpose, // Special Purpose Reimbursable
	"I": ordersmessages.OrdersTypeAccession,      // NAVCAD(Naval Cadet) Accession
	"J": ordersmessages.OrdersTypeAccession,      // Accession Travel Recruits
	"K": ordersmessages.OrdersTypeAccession,      // Non-recruit Accession Travel
	"L": ordersmessages.OrdersTypeTraining,       // Training Travel
	"M": ordersmessages.OrdersTypeRotational,     // Rotational Travel
	"N": ordersmessages.OrdersTypeSeparation,     // Separation Travel
	"O": ordersmessages.OrdersTypeUnitMove,       // Organized Unit/Homeport Change
	"P": ordersmessages.OrdersTypeSeparation,     // Midshipman Separation Travel
	"R": ordersmessages.OrdersTypeOperational,    // Misc. Operational Non-member
	"X": ordersmessages.OrdersTypeEmergencyEvac,  // EMERGENCY NON-MEMBER EVACS
	"Y": ordersmessages.OrdersTypeRotational,     // Misc. Rotational Non-member
	"Z": ordersmessages.OrdersTypeSeparation,     // NAVCAD(Naval Cadet) Separation
}

var officerOrdersTypes = map[string]ordersmessages.OrdersType{
	"0": ordersmessages.OrdersTypeIpcot,          // IPCOT In-place consecutive overseas travel
	"2": ordersmessages.OrdersTypeAccession,      // Accession Travel
	"3": ordersmessages.OrdersTypeTraining,       // Training Travel
	"4": ordersmessages.OrdersTypeOperational,    // Operational Travel
	"5": ordersmessages.OrdersTypeSeparation,     // Separation Travel
	"6": ordersmessages.OrdersTypeUnitMove,       // Organized Unit/Homeport Change
	"7": ordersmessages.OrdersTypeEmergencyEvac,  // Emergency Non-member Evac
	"H": ordersmessages.OrdersTypeSpecialPurpose, // Special Purpose Reimbursable
	"Q": ordersmessages.OrdersTypeRotational,     // Misc. Rotational Non-member
	"S": ordersmessages.OrdersTypeAccession,      // Accession Travel
	"T": ordersmessages.OrdersTypeTraining,       // Training Travel
	"U": ordersmessages.OrdersTypeRotational,     // Rotational Travel
	"V": ordersmessages.OrdersTypeSeparation,     // Separation Travel
	"W": ordersmessages.OrdersTypeUnitMove,       // Organized Unit/Homeport Change
	"X": ordersmessages.OrdersTypeRotational,     // Misc. Rotational Non-member
}
