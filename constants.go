package nursys

// Constants from the API Spec Appendix

// A.2 License types
const (
	LicenseTypeRN   = "RN"   // Registered Nurse
	LicenseTypePN   = "PN"   // Practical Nurse (Vocational Nurse)
	LicenseTypeCNM  = "CNM"  // Certified Nurse Midwife
	LicenseTypeCRNA = "CRNA" // Certified Registered Nurse Anesthetist
	LicenseTypeCNS  = "CNS"  // Clinical Nurse Specialist
	LicenseTypeCNP  = "CNP"  // Certified Nurse Practitioner
)

// A.7 Submission Action Codes
const (
	ActionCodeAdd    = "A" // Either add a new nurse to your nurse list or update an existing nurse
	ActionCodeRemove = "R" // Remove this nurse
)
