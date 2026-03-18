package model

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID                string         `json:"id" gorm:"primaryKey;type:uuid"`
	Name              string         `json:"name" gorm:"size:255;not null"`
	APIKey            string         `json:"api_key" gorm:"size:255;uniqueIndex;not null"`
	ExternalDeviceAPIURL string      `json:"external_device_api_url" gorm:"size:1024"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

type Product struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid"`
	TenantID    string         `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string         `json:"name" gorm:"size:255;not null"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Device struct {
	ID                 string         `json:"id" gorm:"primaryKey;type:uuid"`
	TenantID           string         `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ProductID          string         `json:"product_id" gorm:"type:uuid;not null;index"`
	ExternalDeviceID   string         `json:"external_device_id" gorm:"size:255;not null"`
	CurrentVersion     string         `json:"current_version" gorm:"size:100"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

type Firmware struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid"`
	TenantID    string         `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ProductID   string         `json:"product_id" gorm:"type:uuid;not null;index"`
	Version     string         `json:"version" gorm:"size:100;not null"`
	VersionCode int            `json:"version_code" gorm:"not null"`
	Changelog   string         `json:"changelog" gorm:"type:text"`
	FilePath    string         `json:"-" gorm:"size:1024;not null"`
	FileSize    int64          `json:"file_size" gorm:"not null"`
	MD5         string         `json:"md5" gorm:"size:32;not null"`
	ReleaseNotes string       `json:"release_notes" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type UpgradeType string

const (
	UpgradeTypeSpecified UpgradeType = "specified"
	UpgradeTypeAll       UpgradeType = "all"
	UpgradeTypeGray      UpgradeType = "gray"
)

type TaskStatus string

const (
	TaskStatusCreated   TaskStatus = "created"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusPaused    TaskStatus = "paused"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

type UpgradeTask struct {
	ID                 string         `json:"id" gorm:"primaryKey;type:uuid"`
	TenantID           string         `json:"tenant_id" gorm:"type:uuid;not null;index"`
	ProductID          string         `json:"product_id" gorm:"type:uuid;not null;index"`
	FirmwareID         string         `json:"firmware_id" gorm:"type:uuid;not null"`
	TaskName           string         `json:"task_name" gorm:"size:255;not null"`
	UpgradeType        UpgradeType    `json:"upgrade_type" gorm:"size:20;not null"`
	GrayPercent        int            `json:"gray_percent"`
	TargetDevicesCount int            `json:"target_devices_count"`
	PushRate           int            `json:"push_rate"`
	Status             TaskStatus     `json:"status" gorm:"size:20;not null"`
	CreatedBy          string         `json:"created_by" gorm:"size:255"`
	StartedAt          *time.Time     `json:"started_at"`
	CompletedAt        *time.Time     `json:"completed_at"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

type DeviceUpgradeStatus string

const (
	DeviceStatusPending    DeviceUpgradeStatus = "pending"
	DeviceStatusNotified   DeviceUpgradeStatus = "notified"
	DeviceStatusDownloading DeviceUpgradeStatus = "downloading"
	DeviceStatusDownloaded DeviceUpgradeStatus = "downloaded"
	DeviceStatusInstalling DeviceUpgradeStatus = "installing"
	DeviceStatusSuccess    DeviceUpgradeStatus = "success"
	DeviceStatusFailed     DeviceUpgradeStatus = "failed"
)

type DeviceUpgradeRecord struct {
	ID             string                 `json:"id" gorm:"primaryKey;type:uuid"`
	TaskID         string                 `json:"task_id" gorm:"type:uuid;not null;index"`
	DeviceID       string                 `json:"device_id" gorm:"type:uuid;not null;index"`
	OldVersion     string                 `json:"old_version" gorm:"size:100"`
	NewVersion     string                 `json:"new_version" gorm:"size:100"`
	Status         DeviceUpgradeStatus   `json:"status" gorm:"size:20;not null"`
	ErrorMessage   string                 `json:"error_message" gorm:"type:text"`
	StartedAt      *time.Time            `json:"started_at"`
	FinishedAt     *time.Time            `json:"finished_at"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
}

type ThirdPartyDevice struct {
	DeviceID      string `json:"deviceId"`
	DeviceName    string `json:"deviceName"`
	CurrentVersion string `json:"currentVersion"`
	ProductID     string `json:"productId"`
	Online        bool   `json:"online"`
}

type ThirdPartyDeviceQueryRequest struct {
	ProductID       string   `json:"productId"`
	Percent         int      `json:"percent"`
	ExcludeVersions []string `json:"excludeVersions"`
	OnlyOnline      bool     `json:"onlyOnline"`
	Limit           int      `json:"limit"`
}

type ThirdPartyDeviceQueryResponse struct {
	Code    int `json:"code"`
	Data    struct {
		Total    int                `json:"total"`
		Selected []ThirdPartyDevice `json:"selected"`
	} `json:"data"`
}

type ThirdPartyResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
