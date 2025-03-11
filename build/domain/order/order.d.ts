import _m0 from "protobufjs/minimal";
import { Decimal } from "../decimal/decimal";
import { Denom } from "../denom/denom";
import { MetaData } from "../metadata/metadata";
import { Side } from "../order-properties/order-properties";
export declare const protobufPackage = "order";
/** Type is order type. */
export declare enum OrderType {
    /** ORDER_TYPE_UNSPECIFIED - order_type_unspecified reserves the default value, to protect against unexpected settings. */
    ORDER_TYPE_UNSPECIFIED = 0,
    /** ORDER_TYPE_LIMIT - order_type_limit means that the order is limit order. */
    ORDER_TYPE_LIMIT = 1,
    /** ORDER_TYPE_MARKET - limit order_type_market that the order is market order. */
    ORDER_TYPE_MARKET = 2,
    UNRECOGNIZED = -1
}
export declare function orderTypeFromJSON(object: any): OrderType;
export declare function orderTypeToJSON(object: OrderType): string;
/** TimeInForce is order time in force. */
export declare enum TimeInForce {
    /** TIME_IN_FORCE_UNSPECIFIED - time_in_force_unspecified reserves the default value, to protect against unexpected settings. */
    TIME_IN_FORCE_UNSPECIFIED = 0,
    /** TIME_IN_FORCE_GTC - time_in_force_gtc means that the order remains active until it is fully executed or manually canceled. */
    TIME_IN_FORCE_GTC = 1,
    /**
     * TIME_IN_FORCE_IOC - time_in_force_ioc  means that order must be executed immediately, either in full or partially. Any portion of the
     *  order that cannot be filled immediately is canceled.
     */
    TIME_IN_FORCE_IOC = 2,
    /** TIME_IN_FORCE_FOK - time_in_force_fok means that order must be fully executed or canceled. */
    TIME_IN_FORCE_FOK = 3,
    UNRECOGNIZED = -1
}
export declare function timeInForceFromJSON(object: any): TimeInForce;
export declare function timeInForceToJSON(object: TimeInForce): string;
export declare enum OrderStatus {
    /** ORDER_STATUS_UNSPECIFIED - order_status_unspecified reserves the default value, to protect against unexpected settings. */
    ORDER_STATUS_UNSPECIFIED = 0,
    /** ORDER_STATUS_OPEN - order_status_open means that the order is open with any remaining quantity */
    ORDER_STATUS_OPEN = 1,
    /** ORDER_STATUS_CANCELED - order_status_cancelled means the user has canceled the order. */
    ORDER_STATUS_CANCELED = 2,
    /** ORDER_STATUS_FILLED - order_status_filled means that the order is filled (quantity remaining is 0) */
    ORDER_STATUS_FILLED = 3,
    /** ORDER_STATUS_EXPIRED - order_status_expired means that the order is expired (e.g. a block event has passed the good til block height/time). */
    ORDER_STATUS_EXPIRED = 4,
    UNRECOGNIZED = -1
}
export declare function orderStatusFromJSON(object: any): OrderStatus;
export declare function orderStatusToJSON(object: OrderStatus): string;
/** Unique key is Sequence-Network */
export interface Order {
    /** account is order creator address. */
    Account: string;
    Type: OrderType;
    OrderID: string;
    /** Sequence ID */
    Sequence: number;
    BaseDenom: Denom | undefined;
    QuoteDenom: Denom | undefined;
    /** price is value of one unit of the BaseDenom expressed in terms of the QuoteDenom. */
    Price: number;
    /** quantity is amount of the base BaseDenom being traded. */
    Quantity: Decimal | undefined;
    RemainingQuantity: Decimal | undefined;
    /** Buy or sell */
    Side: Side;
    GoodTil: GoodTil | undefined;
    TimeInForce: TimeInForce;
    /** Time the order was created on chain. This can differ from metadata.CreatedAt which signifies when the record was created in the database */
    BlockTime: Date | undefined;
    /** Maintain the status of the order (tracked for user intent clarification) */
    OrderStatus: OrderStatus;
    OrderFee: number;
    MetaData: MetaData | undefined;
    TXID?: string | undefined;
    BlockHeight: number;
    /** If the order has been enriched with precision data */
    Enriched: boolean;
}
/** GoodTil is a good til order settings. */
export interface GoodTil {
    /** good_til_block_height means that order remains active until a specific blockchain block height is reached. */
    BlockHeight: number;
    /** good_til_block_time means that order remains active until a specific blockchain block time is reached. */
    BlockTime: Date | undefined;
}
export interface Orders {
    Orders: Order[];
    Offset?: number | undefined;
}
export declare const Order: {
    encode(message: Order, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Order;
    fromJSON(object: any): Order;
    toJSON(message: Order): unknown;
    create<I extends {
        Account?: string | undefined;
        Type?: OrderType | undefined;
        OrderID?: string | undefined;
        Sequence?: number | undefined;
        BaseDenom?: {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } | undefined;
        QuoteDenom?: {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } | undefined;
        Price?: number | undefined;
        Quantity?: {
            Value?: number | undefined;
            Exp?: number | undefined;
        } | undefined;
        RemainingQuantity?: {
            Value?: number | undefined;
            Exp?: number | undefined;
        } | undefined;
        Side?: Side | undefined;
        GoodTil?: {
            BlockHeight?: number | undefined;
            BlockTime?: Date | undefined;
        } | undefined;
        TimeInForce?: TimeInForce | undefined;
        BlockTime?: Date | undefined;
        OrderStatus?: OrderStatus | undefined;
        OrderFee?: number | undefined;
        MetaData?: {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } | undefined;
        TXID?: string | undefined;
        BlockHeight?: number | undefined;
        Enriched?: boolean | undefined;
    } & {
        Account?: string | undefined;
        Type?: OrderType | undefined;
        OrderID?: string | undefined;
        Sequence?: number | undefined;
        BaseDenom?: ({
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } & {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } & { [K in Exclude<keyof I["BaseDenom"], keyof Denom>]: never; }) | undefined;
        QuoteDenom?: ({
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } & {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } & { [K_1 in Exclude<keyof I["QuoteDenom"], keyof Denom>]: never; }) | undefined;
        Price?: number | undefined;
        Quantity?: ({
            Value?: number | undefined;
            Exp?: number | undefined;
        } & {
            Value?: number | undefined;
            Exp?: number | undefined;
        } & { [K_2 in Exclude<keyof I["Quantity"], keyof Decimal>]: never; }) | undefined;
        RemainingQuantity?: ({
            Value?: number | undefined;
            Exp?: number | undefined;
        } & {
            Value?: number | undefined;
            Exp?: number | undefined;
        } & { [K_3 in Exclude<keyof I["RemainingQuantity"], keyof Decimal>]: never; }) | undefined;
        Side?: Side | undefined;
        GoodTil?: ({
            BlockHeight?: number | undefined;
            BlockTime?: Date | undefined;
        } & {
            BlockHeight?: number | undefined;
            BlockTime?: Date | undefined;
        } & { [K_4 in Exclude<keyof I["GoodTil"], keyof GoodTil>]: never; }) | undefined;
        TimeInForce?: TimeInForce | undefined;
        BlockTime?: Date | undefined;
        OrderStatus?: OrderStatus | undefined;
        OrderFee?: number | undefined;
        MetaData?: ({
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & { [K_5 in Exclude<keyof I["MetaData"], keyof MetaData>]: never; }) | undefined;
        TXID?: string | undefined;
        BlockHeight?: number | undefined;
        Enriched?: boolean | undefined;
    } & { [K_6 in Exclude<keyof I, keyof Order>]: never; }>(base?: I | undefined): Order;
    fromPartial<I_1 extends {
        Account?: string | undefined;
        Type?: OrderType | undefined;
        OrderID?: string | undefined;
        Sequence?: number | undefined;
        BaseDenom?: {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } | undefined;
        QuoteDenom?: {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } | undefined;
        Price?: number | undefined;
        Quantity?: {
            Value?: number | undefined;
            Exp?: number | undefined;
        } | undefined;
        RemainingQuantity?: {
            Value?: number | undefined;
            Exp?: number | undefined;
        } | undefined;
        Side?: Side | undefined;
        GoodTil?: {
            BlockHeight?: number | undefined;
            BlockTime?: Date | undefined;
        } | undefined;
        TimeInForce?: TimeInForce | undefined;
        BlockTime?: Date | undefined;
        OrderStatus?: OrderStatus | undefined;
        OrderFee?: number | undefined;
        MetaData?: {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } | undefined;
        TXID?: string | undefined;
        BlockHeight?: number | undefined;
        Enriched?: boolean | undefined;
    } & {
        Account?: string | undefined;
        Type?: OrderType | undefined;
        OrderID?: string | undefined;
        Sequence?: number | undefined;
        BaseDenom?: ({
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } & {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } & { [K_7 in Exclude<keyof I_1["BaseDenom"], keyof Denom>]: never; }) | undefined;
        QuoteDenom?: ({
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } & {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } & { [K_8 in Exclude<keyof I_1["QuoteDenom"], keyof Denom>]: never; }) | undefined;
        Price?: number | undefined;
        Quantity?: ({
            Value?: number | undefined;
            Exp?: number | undefined;
        } & {
            Value?: number | undefined;
            Exp?: number | undefined;
        } & { [K_9 in Exclude<keyof I_1["Quantity"], keyof Decimal>]: never; }) | undefined;
        RemainingQuantity?: ({
            Value?: number | undefined;
            Exp?: number | undefined;
        } & {
            Value?: number | undefined;
            Exp?: number | undefined;
        } & { [K_10 in Exclude<keyof I_1["RemainingQuantity"], keyof Decimal>]: never; }) | undefined;
        Side?: Side | undefined;
        GoodTil?: ({
            BlockHeight?: number | undefined;
            BlockTime?: Date | undefined;
        } & {
            BlockHeight?: number | undefined;
            BlockTime?: Date | undefined;
        } & { [K_11 in Exclude<keyof I_1["GoodTil"], keyof GoodTil>]: never; }) | undefined;
        TimeInForce?: TimeInForce | undefined;
        BlockTime?: Date | undefined;
        OrderStatus?: OrderStatus | undefined;
        OrderFee?: number | undefined;
        MetaData?: ({
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & { [K_12 in Exclude<keyof I_1["MetaData"], keyof MetaData>]: never; }) | undefined;
        TXID?: string | undefined;
        BlockHeight?: number | undefined;
        Enriched?: boolean | undefined;
    } & { [K_13 in Exclude<keyof I_1, keyof Order>]: never; }>(object: I_1): Order;
};
export declare const GoodTil: {
    encode(message: GoodTil, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GoodTil;
    fromJSON(object: any): GoodTil;
    toJSON(message: GoodTil): unknown;
    create<I extends {
        BlockHeight?: number | undefined;
        BlockTime?: Date | undefined;
    } & {
        BlockHeight?: number | undefined;
        BlockTime?: Date | undefined;
    } & { [K in Exclude<keyof I, keyof GoodTil>]: never; }>(base?: I | undefined): GoodTil;
    fromPartial<I_1 extends {
        BlockHeight?: number | undefined;
        BlockTime?: Date | undefined;
    } & {
        BlockHeight?: number | undefined;
        BlockTime?: Date | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof GoodTil>]: never; }>(object: I_1): GoodTil;
};
export declare const Orders: {
    encode(message: Orders, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Orders;
    fromJSON(object: any): Orders;
    toJSON(message: Orders): unknown;
    create<I extends {
        Orders?: {
            Account?: string | undefined;
            Type?: OrderType | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            BaseDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            QuoteDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Price?: number | undefined;
            Quantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            RemainingQuantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Side?: Side | undefined;
            GoodTil?: {
                BlockHeight?: number | undefined;
                BlockTime?: Date | undefined;
            } | undefined;
            TimeInForce?: TimeInForce | undefined;
            BlockTime?: Date | undefined;
            OrderStatus?: OrderStatus | undefined;
            OrderFee?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
        }[] | undefined;
        Offset?: number | undefined;
    } & {
        Orders?: ({
            Account?: string | undefined;
            Type?: OrderType | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            BaseDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            QuoteDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Price?: number | undefined;
            Quantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            RemainingQuantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Side?: Side | undefined;
            GoodTil?: {
                BlockHeight?: number | undefined;
                BlockTime?: Date | undefined;
            } | undefined;
            TimeInForce?: TimeInForce | undefined;
            BlockTime?: Date | undefined;
            OrderStatus?: OrderStatus | undefined;
            OrderFee?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
        }[] & ({
            Account?: string | undefined;
            Type?: OrderType | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            BaseDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            QuoteDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Price?: number | undefined;
            Quantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            RemainingQuantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Side?: Side | undefined;
            GoodTil?: {
                BlockHeight?: number | undefined;
                BlockTime?: Date | undefined;
            } | undefined;
            TimeInForce?: TimeInForce | undefined;
            BlockTime?: Date | undefined;
            OrderStatus?: OrderStatus | undefined;
            OrderFee?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
        } & {
            Account?: string | undefined;
            Type?: OrderType | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            BaseDenom?: ({
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } & {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } & { [K in Exclude<keyof I["Orders"][number]["BaseDenom"], keyof Denom>]: never; }) | undefined;
            QuoteDenom?: ({
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } & {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } & { [K_1 in Exclude<keyof I["Orders"][number]["QuoteDenom"], keyof Denom>]: never; }) | undefined;
            Price?: number | undefined;
            Quantity?: ({
                Value?: number | undefined;
                Exp?: number | undefined;
            } & {
                Value?: number | undefined;
                Exp?: number | undefined;
            } & { [K_2 in Exclude<keyof I["Orders"][number]["Quantity"], keyof Decimal>]: never; }) | undefined;
            RemainingQuantity?: ({
                Value?: number | undefined;
                Exp?: number | undefined;
            } & {
                Value?: number | undefined;
                Exp?: number | undefined;
            } & { [K_3 in Exclude<keyof I["Orders"][number]["RemainingQuantity"], keyof Decimal>]: never; }) | undefined;
            Side?: Side | undefined;
            GoodTil?: ({
                BlockHeight?: number | undefined;
                BlockTime?: Date | undefined;
            } & {
                BlockHeight?: number | undefined;
                BlockTime?: Date | undefined;
            } & { [K_4 in Exclude<keyof I["Orders"][number]["GoodTil"], keyof GoodTil>]: never; }) | undefined;
            TimeInForce?: TimeInForce | undefined;
            BlockTime?: Date | undefined;
            OrderStatus?: OrderStatus | undefined;
            OrderFee?: number | undefined;
            MetaData?: ({
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & { [K_5 in Exclude<keyof I["Orders"][number]["MetaData"], keyof MetaData>]: never; }) | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
        } & { [K_6 in Exclude<keyof I["Orders"][number], keyof Order>]: never; })[] & { [K_7 in Exclude<keyof I["Orders"], keyof {
            Account?: string | undefined;
            Type?: OrderType | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            BaseDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            QuoteDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Price?: number | undefined;
            Quantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            RemainingQuantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Side?: Side | undefined;
            GoodTil?: {
                BlockHeight?: number | undefined;
                BlockTime?: Date | undefined;
            } | undefined;
            TimeInForce?: TimeInForce | undefined;
            BlockTime?: Date | undefined;
            OrderStatus?: OrderStatus | undefined;
            OrderFee?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
        }[]>]: never; }) | undefined;
        Offset?: number | undefined;
    } & { [K_8 in Exclude<keyof I, keyof Orders>]: never; }>(base?: I | undefined): Orders;
    fromPartial<I_1 extends {
        Orders?: {
            Account?: string | undefined;
            Type?: OrderType | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            BaseDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            QuoteDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Price?: number | undefined;
            Quantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            RemainingQuantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Side?: Side | undefined;
            GoodTil?: {
                BlockHeight?: number | undefined;
                BlockTime?: Date | undefined;
            } | undefined;
            TimeInForce?: TimeInForce | undefined;
            BlockTime?: Date | undefined;
            OrderStatus?: OrderStatus | undefined;
            OrderFee?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
        }[] | undefined;
        Offset?: number | undefined;
    } & {
        Orders?: ({
            Account?: string | undefined;
            Type?: OrderType | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            BaseDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            QuoteDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Price?: number | undefined;
            Quantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            RemainingQuantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Side?: Side | undefined;
            GoodTil?: {
                BlockHeight?: number | undefined;
                BlockTime?: Date | undefined;
            } | undefined;
            TimeInForce?: TimeInForce | undefined;
            BlockTime?: Date | undefined;
            OrderStatus?: OrderStatus | undefined;
            OrderFee?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
        }[] & ({
            Account?: string | undefined;
            Type?: OrderType | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            BaseDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            QuoteDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Price?: number | undefined;
            Quantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            RemainingQuantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Side?: Side | undefined;
            GoodTil?: {
                BlockHeight?: number | undefined;
                BlockTime?: Date | undefined;
            } | undefined;
            TimeInForce?: TimeInForce | undefined;
            BlockTime?: Date | undefined;
            OrderStatus?: OrderStatus | undefined;
            OrderFee?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
        } & {
            Account?: string | undefined;
            Type?: OrderType | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            BaseDenom?: ({
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } & {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } & { [K_9 in Exclude<keyof I_1["Orders"][number]["BaseDenom"], keyof Denom>]: never; }) | undefined;
            QuoteDenom?: ({
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } & {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } & { [K_10 in Exclude<keyof I_1["Orders"][number]["QuoteDenom"], keyof Denom>]: never; }) | undefined;
            Price?: number | undefined;
            Quantity?: ({
                Value?: number | undefined;
                Exp?: number | undefined;
            } & {
                Value?: number | undefined;
                Exp?: number | undefined;
            } & { [K_11 in Exclude<keyof I_1["Orders"][number]["Quantity"], keyof Decimal>]: never; }) | undefined;
            RemainingQuantity?: ({
                Value?: number | undefined;
                Exp?: number | undefined;
            } & {
                Value?: number | undefined;
                Exp?: number | undefined;
            } & { [K_12 in Exclude<keyof I_1["Orders"][number]["RemainingQuantity"], keyof Decimal>]: never; }) | undefined;
            Side?: Side | undefined;
            GoodTil?: ({
                BlockHeight?: number | undefined;
                BlockTime?: Date | undefined;
            } & {
                BlockHeight?: number | undefined;
                BlockTime?: Date | undefined;
            } & { [K_13 in Exclude<keyof I_1["Orders"][number]["GoodTil"], keyof GoodTil>]: never; }) | undefined;
            TimeInForce?: TimeInForce | undefined;
            BlockTime?: Date | undefined;
            OrderStatus?: OrderStatus | undefined;
            OrderFee?: number | undefined;
            MetaData?: ({
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & { [K_14 in Exclude<keyof I_1["Orders"][number]["MetaData"], keyof MetaData>]: never; }) | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
        } & { [K_15 in Exclude<keyof I_1["Orders"][number], keyof Order>]: never; })[] & { [K_16 in Exclude<keyof I_1["Orders"], keyof {
            Account?: string | undefined;
            Type?: OrderType | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            BaseDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            QuoteDenom?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Price?: number | undefined;
            Quantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            RemainingQuantity?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Side?: Side | undefined;
            GoodTil?: {
                BlockHeight?: number | undefined;
                BlockTime?: Date | undefined;
            } | undefined;
            TimeInForce?: TimeInForce | undefined;
            BlockTime?: Date | undefined;
            OrderStatus?: OrderStatus | undefined;
            OrderFee?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
        }[]>]: never; }) | undefined;
        Offset?: number | undefined;
    } & { [K_17 in Exclude<keyof I_1, keyof Orders>]: never; }>(object: I_1): Orders;
};
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P : P & {
    [K in keyof P]: Exact<P[K], I[K]>;
} & {
    [K in Exclude<keyof I, KeysOfUnion<P>>]: never;
};
export {};
