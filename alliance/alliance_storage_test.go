package alliance

import (
	"fmt"
	"testing"
)

func TestAllianceStorage_AllianceStorageClearup(t *testing.T) {
	type fields struct {
		AllianceId  string
		OwnerId     string
		Items       []ItemDesc
		MaxCapacity int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				Items: []ItemDesc{
					{ItemType: 1, Number: 1},
					{ItemType: 1, Number: 1},
					{ItemType: 1, Number: 3},
					{ItemType: 3, Number: 1},
					{ItemType: 3, Number: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "overflow",
			fields: fields{
				Items: []ItemDesc{
					{ItemType: 1, Number: 4},
					{ItemType: 1, Number: 4},
					{ItemType: 1, Number: 3},
					{ItemType: 3, Number: 1},
					{ItemType: 3, Number: 1},
					{ItemType: 1, Number: 3},
				},
			},
			wantErr: false,
		},
		{
			name: "empty hole",
			fields: fields{
				Items: []ItemDesc{
					{ItemType: 1, Number: 1},
					{ItemType: 5, Number: 1},
					{ItemType: EmptyGridItemType, Number: 0},
					{ItemType: 3, Number: 1},
					{ItemType: 3, Number: 1},
					{ItemType: 5, Number: 1},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := &Alliance{
				AllianceId:  tt.fields.AllianceId,
				OwnerId:     tt.fields.OwnerId,
				Items:       tt.fields.Items,
				MaxCapacity: tt.fields.MaxCapacity,
			}
			if err := as.clearup(); (err != nil) != tt.wantErr {
				t.Errorf("AllianceStorage.AllianceStorageClearup() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(as)
		})
	}
}

func TestAllianceStorage_add(t *testing.T) {
	type fields struct {
		AllianceId  string
		OwnerId     string
		Items       []ItemDesc
		MaxCapacity int
	}
	type args struct {
		index    int
		itemtype int32
		num      int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "demo1",
			fields: fields{
				Items: []ItemDesc{
					{ItemType: 1, Number: 1},
					{ItemType: 1, Number: 1},
					{ItemType: 1, Number: 3},
					{ItemType: 3, Number: 1},
					{ItemType: 3, Number: 1},
				},
			},
			args:    args{index: 1, itemtype: 1, num: 3},
			wantErr: false,
		},
		{
			name: "overflow",
			fields: fields{
				Items: []ItemDesc{
					{ItemType: 1, Number: 1},
					{ItemType: 1, Number: 1},
					{ItemType: EmptyGridItemType, Number: 0},
					{ItemType: 3, Number: 1},
					{ItemType: 3, Number: 1},
					{ItemType: EmptyGridItemType, Number: 0},
					{ItemType: EmptyGridItemType, Number: 0},
					{ItemType: EmptyGridItemType, Number: 0},
				},
			},
			args:    args{index: 2, itemtype: 1, num: 8},
			wantErr: false,
		},
		{
			name: "full",
			fields: fields{
				Items: []ItemDesc{
					{ItemType: 1, Number: 1},
					{ItemType: 1, Number: 1},
					{ItemType: EmptyGridItemType, Number: 0},
					{ItemType: 3, Number: 1},
					{ItemType: 3, Number: 1},
				},
			},
			args:    args{index: 2, itemtype: 1, num: 8},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := &Alliance{
				AllianceId:  tt.fields.AllianceId,
				OwnerId:     tt.fields.OwnerId,
				Items:       tt.fields.Items,
				MaxCapacity: tt.fields.MaxCapacity,
			}
			if err := as.add(tt.args.index, tt.args.itemtype, tt.args.num); (err != nil) != tt.wantErr {
				t.Errorf("AllianceStorage.add() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(as)
		})
	}
}
