package data_saver

// DataSaver HardSave -- saves and flushes Save -- only saves
type DataSaver interface {
	HardSave([]byte) (string, error)
}
