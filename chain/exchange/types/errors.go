package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrOrderInvalid                       = sdkerrors.Register(ModuleName, 1, "failed to validate order")
	ErrSpotMarketNotFound                 = sdkerrors.Register(ModuleName, 2, "spot market not found")
	ErrSpotMarketExists                   = sdkerrors.Register(ModuleName, 3, "spot market exists")
	ErrBadField                           = sdkerrors.Register(ModuleName, 4, "struct field error")
	ErrMarketInvalid                      = sdkerrors.Register(ModuleName, 5, "failed to validate derivative market")
	ErrInsufficientDeposit                = sdkerrors.Register(ModuleName, 6, "subaccount has insufficient deposits")
	ErrUnrecognizedOrderType              = sdkerrors.Register(ModuleName, 7, "unrecognized order type")
	ErrInsufficientPositionQuantity       = sdkerrors.Register(ModuleName, 8, "position quantity insufficient for order")
	ErrOrderHashInvalid                   = sdkerrors.Register(ModuleName, 9, "order hash is not valid")
	ErrBadSubaccountID                    = sdkerrors.Register(ModuleName, 10, "subaccount id is not valid")
	ErrInvalidTicker                      = sdkerrors.Register(ModuleName, 11, "invalid ticker")
	ErrInvalidBaseDenom                   = sdkerrors.Register(ModuleName, 12, "invalid base denom")
	ErrInvalidQuoteDenom                  = sdkerrors.Register(ModuleName, 13, "invalid quote denom")
	ErrInvalidOracle                      = sdkerrors.Register(ModuleName, 14, "invalid oracle")
	ErrInvalidExpiry                      = sdkerrors.Register(ModuleName, 15, "invalid expiry")
	ErrInvalidPrice                       = sdkerrors.Register(ModuleName, 16, "invalid price")
	ErrInvalidQuantity                    = sdkerrors.Register(ModuleName, 17, "invalid quantity")
	ErrUnsupportedOracleType              = sdkerrors.Register(ModuleName, 18, "unsupported oracle type")
	ErrOrderDoesntExist                   = sdkerrors.Register(ModuleName, 19, "order doesnt exist")
	ErrOrderbookFillInvalid               = sdkerrors.Register(ModuleName, 20, "spot limit orderbook fill invalid")
	ErrPerpetualMarketExists              = sdkerrors.Register(ModuleName, 21, "perpetual market exists")
	ErrExpiryFuturesMarketExists          = sdkerrors.Register(ModuleName, 22, "expiry futures market exists")
	ErrExpiryFuturesMarketExpired         = sdkerrors.Register(ModuleName, 23, "expiry futures market expired")
	ErrNoLiquidity                        = sdkerrors.Register(ModuleName, 24, "no liquidity on the orderbook!")
	ErrSlippageExceedsWorstPrice          = sdkerrors.Register(ModuleName, 25, "Orderbook liquidity cannot satisfy current worst price")
	ErrInsufficientOrderMargin            = sdkerrors.Register(ModuleName, 26, "Order has insufficient margin")
	ErrDerivativeMarketNotFound           = sdkerrors.Register(ModuleName, 27, "Derivative market not found")
	ErrPositionNotFound                   = sdkerrors.Register(ModuleName, 28, "Position not found")
	ErrInvalidReduceOnlyPositionDirection = sdkerrors.Register(ModuleName, 29, "Position direction does not oppose the reduce-only order")
	ErrPriceSurpassesBankruptcyPrice      = sdkerrors.Register(ModuleName, 30, "Price Surpasses Bankruptcy Price")
	ErrPositionNotLiquidable              = sdkerrors.Register(ModuleName, 31, "Position not liquidable")
	ErrInvalidTriggerPrice                = sdkerrors.Register(ModuleName, 32, "invalid trigger price")
	ErrInvalidOracleType                  = sdkerrors.Register(ModuleName, 33, "invalid oracle type")
	ErrInvalidPriceTickSize               = sdkerrors.Register(ModuleName, 34, "invalid minimum price tick size")
	ErrInvalidQuantityTickSize            = sdkerrors.Register(ModuleName, 35, "invalid minimum quantity tick size")
	ErrInvalidMargin                      = sdkerrors.Register(ModuleName, 36, "invalid minimum order margin")
	ErrExceedsOrderSideCount              = sdkerrors.Register(ModuleName, 37, "Exceeds order side count")
	ErrDerivativeMarketOrderAlreadyExists = sdkerrors.Register(ModuleName, 38, "Subaccount cannot place a market order when a market order or limit order in the same market was already placed in same block")
	ErrDerivativeLimitOrderAlreadyExists  = sdkerrors.Register(ModuleName, 39, "Subaccount cannot place a limit order when a market order in the same market was already placed in the same block")
	ErrMarketLaunchProposalAlreadyExists  = sdkerrors.Register(ModuleName, 40, "An equivalent market launch proposal already exists.")
	ErrInvalidMarketStatus                = sdkerrors.Register(ModuleName, 41, "Invalid Market Status")
	ErrSameDenoms                         = sdkerrors.Register(ModuleName, 42, "base denom cannot be same with quote denom")
	ErrSameOracles                        = sdkerrors.Register(ModuleName, 43, "oracle base cannot be same with oracle quote")
	ErrFeeRatesRelation                   = sdkerrors.Register(ModuleName, 44, "MakerFeeRate cannot be greater than TakerFeeRate")
	ErrMarginsRelation                    = sdkerrors.Register(ModuleName, 45, "MaintenanceMarginRatio cannot be greater than InitialMarginRatio")
	ErrExceedsMaxOracleScaleFactor        = sdkerrors.Register(ModuleName, 46, "OracleScaleFactor cannot be greater than MaxOracleScaleFactor")
	ErrSpotExchangeNotEnabled             = sdkerrors.Register(ModuleName, 47, "Spot exchange is not enabled yet")
	ErrDerivativesExchangeNotEnabled      = sdkerrors.Register(ModuleName, 48, "Derivatives exchange is not enabled yet")
)
