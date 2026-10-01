package models

import "time"

type Subscription struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SubscriberID uint      `gorm:"uniqueIndex:idx_sub_chan" json:"subscriber_id"`
	ChannelID    uint      `gorm:"uniqueIndex:idx_sub_chan" json:"channel_id"`
	CreatedAt    time.Time `json:"created_at"`

	Subscriber   User      `gorm:"foreignKey:SubscriberID" json:"subscriber,omitempty"`
	Channel      User      `gorm:"foreignKey:ChannelID" json:"channel,omitempty"`
}
