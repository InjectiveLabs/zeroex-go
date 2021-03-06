package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/InjectiveLabs/sdk-go/chain/oracle/types"
)

const RouterKey = ModuleName

var (
	_ sdk.Msg = &MsgDeposit{}
	_ sdk.Msg = &MsgWithdraw{}
	_ sdk.Msg = &MsgCreateSpotLimitOrder{}
	_ sdk.Msg = &MsgBatchCreateSpotLimitOrders{}
	_ sdk.Msg = &MsgCreateSpotMarketOrder{}
	_ sdk.Msg = &MsgCancelSpotOrder{}
	_ sdk.Msg = &MsgBatchCancelSpotOrders{}
	_ sdk.Msg = &MsgCreateDerivativeLimitOrder{}
	_ sdk.Msg = &MsgBatchCreateDerivativeLimitOrders{}
	_ sdk.Msg = &MsgCreateDerivativeMarketOrder{}
	_ sdk.Msg = &MsgCancelDerivativeOrder{}
	_ sdk.Msg = &MsgBatchCancelDerivativeOrders{}
	_ sdk.Msg = &MsgSubaccountTransfer{}
	_ sdk.Msg = &MsgExternalTransfer{}
	_ sdk.Msg = &MsgIncreasePositionMargin{}
	_ sdk.Msg = &MsgLiquidatePosition{}
	_ sdk.Msg = &MsgInstantSpotMarketLaunch{}
	_ sdk.Msg = &MsgInstantPerpetualMarketLaunch{}
	_ sdk.Msg = &MsgInstantExpiryFuturesMarketLaunch{}
)

func (o *SpotOrder) ValidateBasic(senderAddr sdk.AccAddress) error {
	if o.MarketId == "" {
		return sdkerrors.Wrap(ErrMarketInvalid, o.MarketId)
	}
	switch o.OrderType {
	case OrderType_BUY, OrderType_SELL:
		// do nothing
	default:
		return sdkerrors.Wrap(ErrUnrecognizedOrderType, string(o.OrderType))
	}
	if o.TriggerPrice != nil && (o.TriggerPrice.IsNil() || o.TriggerPrice.LT(sdk.ZeroDec()) || o.TriggerPrice.GT(MaxOrderPrice)) {
		return sdkerrors.Wrap(ErrInvalidTriggerPrice, o.TriggerPrice.String())
	}

	_, err := sdk.AccAddressFromBech32(o.OrderInfo.FeeRecipient)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, o.OrderInfo.FeeRecipient)
	}

	return o.OrderInfo.ValidateBasic(senderAddr)
}

func (o *OrderInfo) ValidateBasic(senderAddr sdk.AccAddress) error {
	subaccountAddress, ok := IsValidSubaccountID(o.SubaccountId)
	if !ok {
		return sdkerrors.Wrap(ErrBadSubaccountID, o.SubaccountId)
	}
	if !bytes.Equal(subaccountAddress.Bytes(), senderAddr.Bytes()) {
		return sdkerrors.Wrap(ErrBadSubaccountID, senderAddr.String())
	}

	if o.Quantity.IsNil() || o.Quantity.LTE(sdk.ZeroDec()) || o.Quantity.GT(MaxOrderQuantity) {
		return sdkerrors.Wrap(ErrInvalidQuantity, o.Quantity.String())
	}

	if o.Price.IsNil() || o.Price.LTE(sdk.ZeroDec()) || o.Price.GT(MaxOrderPrice) {
		return sdkerrors.Wrap(ErrInvalidPrice, o.Price.String())
	}
	return nil
}

func (o *DerivativeOrder) ValidateBasic(senderAddr sdk.AccAddress) error {
	if o.MarketId == "" {
		return sdkerrors.Wrap(ErrMarketInvalid, o.MarketId)
	}
	switch o.OrderType {
	case OrderType_BUY, OrderType_SELL:
		// do nothing
	default:
		return sdkerrors.Wrap(ErrUnrecognizedOrderType, string(o.OrderType))
	}

	if o.Margin.IsNil() || o.Margin.LT(sdk.ZeroDec()) || o.Margin.GT(MaxOrderPrice) {
		return sdkerrors.Wrap(ErrInsufficientOrderMargin, o.Margin.String())
	}
	if o.TriggerPrice != nil && (o.TriggerPrice.IsNil() || o.TriggerPrice.LT(sdk.ZeroDec()) || o.TriggerPrice.GT(MaxOrderPrice)) {
		return sdkerrors.Wrap(ErrInvalidTriggerPrice, o.TriggerPrice.String())
	}

	_, err := sdk.AccAddressFromBech32(o.OrderInfo.FeeRecipient)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, o.OrderInfo.FeeRecipient)
	}

	return o.OrderInfo.ValidateBasic(senderAddr)
}

func (o *OrderData) ValidateBasic(senderAddr sdk.AccAddress) error {
	if o.MarketId == "" {
		return sdkerrors.Wrap(ErrMarketInvalid, o.MarketId)
	}

	subaccountAddress, ok := IsValidSubaccountID(o.SubaccountId)
	if !ok {
		return sdkerrors.Wrap(ErrBadSubaccountID, o.SubaccountId)
	}

	if !bytes.Equal(subaccountAddress.Bytes(), senderAddr.Bytes()) {
		return sdkerrors.Wrap(ErrBadSubaccountID, senderAddr.String())
	}

	if ok = IsValidOrderHash(o.OrderHash); !ok {
		return sdkerrors.Wrap(ErrOrderHashInvalid, o.OrderHash)
	}

	return nil
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgDeposit) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgDeposit) Type() string { return "msgDeposit" }

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg MsgDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)

	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if len(msg.SubaccountId) == 0 {
		return nil
	} else if _, ok := IsValidSubaccountID(msg.SubaccountId); !ok {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.SubaccountId)
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgDeposit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgWithdraw) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgWithdraw) Type() string { return "msgWithdraw" }

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg MsgWithdraw) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	subaccountAddress, ok := IsValidSubaccountID(msg.SubaccountId)
	if !ok {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.SubaccountId)
	}
	if !bytes.Equal(subaccountAddress.Bytes(), senderAddr.Bytes()) {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.Sender)
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgWithdraw) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgWithdraw) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgInstantSpotMarketLaunch) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgInstantSpotMarketLaunch) Type() string { return "instantSpotMarketLaunch" }

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg MsgInstantSpotMarketLaunch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	if msg.Ticker == "" {
		return sdkerrors.Wrap(ErrInvalidTicker, "ticker should not be empty")
	}
	if msg.BaseDenom == "" {
		return sdkerrors.Wrap(ErrInvalidBaseDenom, "base denom should not be empty")
	}
	if msg.QuoteDenom == "" {
		return sdkerrors.Wrap(ErrInvalidQuoteDenom, "quote denom should not be empty")
	}
	if msg.BaseDenom == msg.QuoteDenom {
		return ErrSameDenoms
	}

	if err := ValidateTickSize(msg.MinPriceTickSize); err != nil {
		return sdkerrors.Wrap(ErrInvalidPriceTickSize, err.Error())
	}
	if err := ValidateTickSize(msg.MinQuantityTickSize); err != nil {
		return sdkerrors.Wrap(ErrInvalidQuantityTickSize, err.Error())
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgInstantSpotMarketLaunch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgInstantSpotMarketLaunch) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgInstantPerpetualMarketLaunch) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgInstantPerpetualMarketLaunch) Type() string { return "instantPerpetualMarketLaunch" }

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg MsgInstantPerpetualMarketLaunch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	if msg.Ticker == "" {
		return sdkerrors.Wrap(ErrInvalidTicker, "ticker should not be empty")
	}
	if msg.QuoteDenom == "" {
		return sdkerrors.Wrap(ErrInvalidQuoteDenom, "quote denom should not be empty")
	}
	if msg.OracleBase == "" {
		return sdkerrors.Wrap(ErrInvalidOracle, "oracle base should not be empty")
	}
	if msg.OracleQuote == "" {
		return sdkerrors.Wrap(ErrInvalidOracle, "oracle quote should not be empty")
	}
	if msg.OracleBase == msg.OracleQuote {
		return ErrSameOracles
	}
	switch msg.OracleType {
	case types.OracleType_Band, types.OracleType_PriceFeed, types.OracleType_Coinbase, types.OracleType_Chainlink, types.OracleType_Razor,
		types.OracleType_Dia, types.OracleType_API3, types.OracleType_Uma, types.OracleType_Pyth, types.OracleType_BandIBC:

	default:
		return sdkerrors.Wrap(ErrInvalidOracleType, msg.OracleType.String())
	}
	if msg.OracleScaleFactor > MaxOracleScaleFactor {
		return ErrExceedsMaxOracleScaleFactor
	}
	if err := ValidateFee(msg.MakerFeeRate); err != nil {
		return err
	}
	if err := ValidateFee(msg.TakerFeeRate); err != nil {
		return err
	}
	if err := ValidateMarginRatio(msg.InitialMarginRatio); err != nil {
		return err
	}
	if err := ValidateMarginRatio(msg.MaintenanceMarginRatio); err != nil {
		return err
	}
	if msg.MakerFeeRate.GT(msg.TakerFeeRate) {
		return ErrFeeRatesRelation
	}
	if msg.InitialMarginRatio.LT(msg.MaintenanceMarginRatio) {
		return ErrMarginsRelation
	}
	if err := ValidateTickSize(msg.MinPriceTickSize); err != nil {
		return sdkerrors.Wrap(ErrInvalidPriceTickSize, err.Error())
	}
	if err := ValidateTickSize(msg.MinQuantityTickSize); err != nil {
		return sdkerrors.Wrap(ErrInvalidQuantityTickSize, err.Error())
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgInstantPerpetualMarketLaunch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgInstantPerpetualMarketLaunch) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgInstantExpiryFuturesMarketLaunch) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgInstantExpiryFuturesMarketLaunch) Type() string {
	return "instantExpiryFuturesMarketLaunch"
}

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg MsgInstantExpiryFuturesMarketLaunch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	if msg.Ticker == "" {
		return sdkerrors.Wrap(ErrInvalidTicker, "ticker should not be empty")
	}
	if msg.QuoteDenom == "" {
		return sdkerrors.Wrap(ErrInvalidQuoteDenom, "quote denom should not be empty")
	}
	if msg.OracleBase == "" {
		return sdkerrors.Wrap(ErrInvalidOracle, "oracle base should not be empty")
	}
	if msg.OracleQuote == "" {
		return sdkerrors.Wrap(ErrInvalidOracle, "oracle quote should not be empty")
	}
	if msg.OracleBase == msg.OracleQuote {
		return ErrSameOracles
	}
	switch msg.OracleType {
	case types.OracleType_Band, types.OracleType_PriceFeed, types.OracleType_Coinbase, types.OracleType_Chainlink, types.OracleType_Razor,
		types.OracleType_Dia, types.OracleType_API3, types.OracleType_Uma, types.OracleType_Pyth, types.OracleType_BandIBC:

	default:
		return sdkerrors.Wrap(ErrInvalidOracleType, msg.OracleType.String())
	}
	if msg.OracleScaleFactor > MaxOracleScaleFactor {
		return ErrExceedsMaxOracleScaleFactor
	}
	if msg.Expiry <= 0 {
		return sdkerrors.Wrap(ErrInvalidExpiry, "expiry should not be empty")
	}
	if err := ValidateFee(msg.MakerFeeRate); err != nil {
		return err
	}
	if err := ValidateFee(msg.TakerFeeRate); err != nil {
		return err
	}
	if err := ValidateMarginRatio(msg.InitialMarginRatio); err != nil {
		return err
	}
	if err := ValidateMarginRatio(msg.MaintenanceMarginRatio); err != nil {
		return err
	}
	if msg.MakerFeeRate.GT(msg.TakerFeeRate) {
		return ErrFeeRatesRelation
	}
	if msg.InitialMarginRatio.LT(msg.MaintenanceMarginRatio) {
		return ErrMarginsRelation
	}
	if err := ValidateTickSize(msg.MinPriceTickSize); err != nil {
		return sdkerrors.Wrap(ErrInvalidPriceTickSize, err.Error())
	}
	if err := ValidateTickSize(msg.MinQuantityTickSize); err != nil {
		return sdkerrors.Wrap(ErrInvalidQuantityTickSize, err.Error())
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgInstantExpiryFuturesMarketLaunch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgInstantExpiryFuturesMarketLaunch) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgCreateSpotLimitOrder) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgCreateSpotLimitOrder) Type() string { return "createSpotLimitOrder" }

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg MsgCreateSpotLimitOrder) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // We don't need to check if sender is empty.
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	if err := msg.Order.ValidateBasic(senderAddr); err != nil {
		return err
	}
	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgCreateSpotLimitOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgCreateSpotLimitOrder) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgBatchCreateSpotLimitOrders) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgBatchCreateSpotLimitOrders) Type() string { return "batchCreateSpotLimitOrders" }

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg MsgBatchCreateSpotLimitOrders) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // We don't need to check if sender is empty.
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	for idx := range msg.Orders {
		order := msg.Orders[idx]
		if err := order.ValidateBasic(senderAddr); err != nil {
			return err
		}
	}
	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgBatchCreateSpotLimitOrders) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgBatchCreateSpotLimitOrders) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg MsgCreateSpotMarketOrder) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg MsgCreateSpotMarketOrder) Type() string { return "createSpotMarketOrder" }

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg MsgCreateSpotMarketOrder) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	if err := msg.Order.ValidateBasic(senderAddr); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgCreateSpotMarketOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg MsgCreateSpotMarketOrder) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg *MsgCancelSpotOrder) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg *MsgCancelSpotOrder) Type() string { return "cancelSpotOrder" }

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg *MsgCancelSpotOrder) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	orderData := OrderData{
		MarketId:     msg.MarketId,
		SubaccountId: msg.SubaccountId,
		OrderHash:    msg.OrderHash,
	}
	return orderData.ValidateBasic(senderAddr)
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgCancelSpotOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg *MsgCancelSpotOrder) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg *MsgBatchCancelSpotOrders) Route() string { return RouterKey }

// Type implements the sdk.Msg interface. It should return the action.
func (msg *MsgBatchCancelSpotOrders) Type() string { return "batchCancelSpotOrders" }

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg *MsgBatchCancelSpotOrders) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	for idx := range msg.Data {
		if err := msg.Data[idx].ValidateBasic(senderAddr); err != nil {
			return err
		}
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgBatchCancelSpotOrders) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg *MsgBatchCancelSpotOrders) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route should return the name of the module
func (msg MsgCreateDerivativeLimitOrder) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateDerivativeLimitOrder) Type() string { return "createDerivativeLimitOrder" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateDerivativeLimitOrder) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	if err := msg.Order.ValidateBasic(senderAddr); err != nil {
		return err
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCreateDerivativeLimitOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateDerivativeLimitOrder) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route should return the name of the module
func (msg MsgBatchCreateDerivativeLimitOrders) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBatchCreateDerivativeLimitOrders) Type() string {
	return "batchCreateDerivativeLimitOrder"
}

// ValidateBasic runs stateless checks on the message
func (msg MsgBatchCreateDerivativeLimitOrders) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	for idx := range msg.Orders {
		order := msg.Orders[idx]
		if err := order.ValidateBasic(senderAddr); err != nil {
			return err
		}
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgBatchCreateDerivativeLimitOrders) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBatchCreateDerivativeLimitOrders) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route should return the name of the module
func (msg MsgCreateDerivativeMarketOrder) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateDerivativeMarketOrder) Type() string { return "createDerivativeMarketOrder" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateDerivativeMarketOrder) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	if err := msg.Order.ValidateBasic(senderAddr); err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCreateDerivativeMarketOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateDerivativeMarketOrder) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg *MsgCancelDerivativeOrder) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface. It should return the action.
func (msg *MsgCancelDerivativeOrder) Type() string {
	return "cancelDerivativeOrder"
}

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg *MsgCancelDerivativeOrder) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	orderData := OrderData{
		MarketId:     msg.MarketId,
		SubaccountId: msg.SubaccountId,
		OrderHash:    msg.OrderHash,
	}
	return orderData.ValidateBasic(senderAddr)
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgCancelDerivativeOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg *MsgCancelDerivativeOrder) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// Route implements the sdk.Msg interface. It should return the name of the module
func (msg *MsgBatchCancelDerivativeOrders) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface. It should return the action.
func (msg *MsgBatchCancelDerivativeOrders) Type() string {
	return "batchCancelDerivativeOrder"
}

// ValidateBasic implements the sdk.Msg interface. It runs stateless checks on the message
func (msg *MsgBatchCancelDerivativeOrders) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	for idx := range msg.Data {
		if err := msg.Data[idx].ValidateBasic(senderAddr); err != nil {
			return err
		}
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface. It encodes the message for signing
func (msg *MsgBatchCancelDerivativeOrders) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It defines whose signature is required
func (msg *MsgBatchCancelDerivativeOrders) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgSubaccountTransfer) Route() string {
	return RouterKey
}

func (msg *MsgSubaccountTransfer) Type() string {
	return "subaccountTransfer"
}

func (msg *MsgSubaccountTransfer) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	subaccountAddress, ok := IsValidSubaccountID(msg.SourceSubaccountId)
	if !ok {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.SourceSubaccountId)
	}
	destSubaccountAddress, ok := IsValidSubaccountID(msg.DestinationSubaccountId)
	if !ok {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.DestinationSubaccountId)
	}
	if !bytes.Equal(subaccountAddress.Bytes(), destSubaccountAddress.Bytes()) {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.DestinationSubaccountId)
	}
	if !bytes.Equal(subaccountAddress.Bytes(), senderAddr.Bytes()) {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.Sender)
	}
	return nil
}

func (msg *MsgSubaccountTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgSubaccountTransfer) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgExternalTransfer) Route() string {
	return RouterKey
}

func (msg *MsgExternalTransfer) Type() string {
	return "externalTransfer"
}

func (msg *MsgExternalTransfer) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	sourceSubaccountAddress, ok := IsValidSubaccountID(msg.SourceSubaccountId)
	if !ok {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.SourceSubaccountId)
	}

	_, ok = IsValidSubaccountID(msg.DestinationSubaccountId)
	if !ok {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.DestinationSubaccountId)
	}
	if !bytes.Equal(sourceSubaccountAddress.Bytes(), senderAddr.Bytes()) {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.Sender)
	}
	return nil
}

func (msg *MsgExternalTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgExternalTransfer) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgIncreasePositionMargin) Route() string {
	return RouterKey
}

func (msg *MsgIncreasePositionMargin) Type() string {
	return "increasePositionMargin"
}

func (msg *MsgIncreasePositionMargin) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	if msg.MarketId == "" {
		return sdkerrors.Wrap(ErrMarketInvalid, msg.MarketId)
	}
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	sourceSubaccountAddress, ok := IsValidSubaccountID(msg.SourceSubaccountId)
	if !ok {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.SourceSubaccountId)
	}
	if !bytes.Equal(sourceSubaccountAddress.Bytes(), senderAddr.Bytes()) {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.Sender)
	}

	_, ok = IsValidSubaccountID(msg.DestinationSubaccountId)
	if !ok {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.DestinationSubaccountId)
	}

	return nil
}

func (msg *MsgIncreasePositionMargin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgIncreasePositionMargin) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgLiquidatePosition) Route() string {
	return RouterKey
}

func (msg *MsgLiquidatePosition) Type() string {
	return "liquidatePosition"
}

func (msg *MsgLiquidatePosition) ValidateBasic() error {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)

	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	if msg.MarketId == "" {
		return sdkerrors.Wrap(ErrMarketInvalid, msg.MarketId)
	}

	_, ok := IsValidSubaccountID(msg.SubaccountId)
	if !ok {
		return sdkerrors.Wrap(ErrBadSubaccountID, msg.SubaccountId)
	}

	if msg.Order != nil {
		if err := msg.Order.ValidateBasic(senderAddr); err != nil {
			return err
		}
	}

	return nil
}

func (msg *MsgLiquidatePosition) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgLiquidatePosition) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
