import _m0 from "protobufjs/minimal";
import { Decimal } from "../decimal/decimal";
import { Denom } from "../denom/denom";
import { MetaData } from "../metadata/metadata";
import { Side } from "../order-properties/order-properties";
export declare const protobufPackage = "trade";
/** Key in store is TXID-Sequence-Metadata.Network */
export interface Trade {
    Account: string;
    /** User assigned order reference */
    OrderID: string;
    /** The sequence number of the order, assigned by the DEX (guaranteed unique value for the order) */
    Sequence: number;
    Amount: Decimal | undefined;
    Price: number;
    Denom1: Denom | undefined;
    Denom2: Denom | undefined;
    /** The buy/sell (e.g. did the user place a buy or sell order) */
    Side: Side;
    /** The time the trade was executed in UTC */
    BlockTime: Date | undefined;
    /** Standard storage related fields */
    MetaData: MetaData | undefined;
    TXID?: string | undefined;
    BlockHeight: number;
    /** If the trade has been enriched with precision data */
    Enriched: boolean;
    /** USD representation of the trade values and trading fee (fixed base for easy data comparisson in reports etc) */
    USD?: number | undefined;
    /**
     * Trades get stored in alphabetical order of the denom pair.
     * Data is "uninverted" on retrieval and
     * this flag only indicates that the denoms as seen in the record are not in the original order
     */
    Inverted: boolean;
}
export interface Trades {
    Trades: Trade[];
}
export interface TradePair {
    Denom1: Denom | undefined;
    Denom2: Denom | undefined;
    MetaData: MetaData | undefined;
    PriceTick?: Decimal | undefined;
    QuantityStep?: number | undefined;
}
export interface TradePairs {
    TradePairs: TradePair[];
    Offset?: number | undefined;
}
export declare const Trade: {
    encode(message: Trade, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Trade;
    fromJSON(object: any): Trade;
    toJSON(message: Trade): unknown;
    create<I extends {
        Account?: string | undefined;
        OrderID?: string | undefined;
        Sequence?: number | undefined;
        Amount?: {
            Value?: number | undefined;
            Exp?: number | undefined;
        } | undefined;
        Price?: number | undefined;
        Denom1?: {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } | undefined;
        Denom2?: {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } | undefined;
        Side?: Side | undefined;
        BlockTime?: Date | undefined;
        MetaData?: {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } | undefined;
        TXID?: string | undefined;
        BlockHeight?: number | undefined;
        Enriched?: boolean | undefined;
        USD?: number | undefined;
        Inverted?: boolean | undefined;
    } & {
        Account?: string | undefined;
        OrderID?: string | undefined;
        Sequence?: number | undefined;
        Amount?: ({
            Value?: number | undefined;
            Exp?: number | undefined;
        } & {
            Value?: number | undefined;
            Exp?: number | undefined;
        } & { [K in Exclude<keyof I["Amount"], keyof Decimal>]: never; }) | undefined;
        Price?: number | undefined;
        Denom1?: ({
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
        } & { [K_1 in Exclude<keyof I["Denom1"], keyof Denom>]: never; }) | undefined;
        Denom2?: ({
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
        } & { [K_2 in Exclude<keyof I["Denom2"], keyof Denom>]: never; }) | undefined;
        Side?: Side | undefined;
        BlockTime?: Date | undefined;
        MetaData?: ({
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & { [K_3 in Exclude<keyof I["MetaData"], keyof MetaData>]: never; }) | undefined;
        TXID?: string | undefined;
        BlockHeight?: number | undefined;
        Enriched?: boolean | undefined;
        USD?: number | undefined;
        Inverted?: boolean | undefined;
    } & { [K_4 in Exclude<keyof I, keyof Trade>]: never; }>(base?: I | undefined): Trade;
    fromPartial<I_1 extends {
        Account?: string | undefined;
        OrderID?: string | undefined;
        Sequence?: number | undefined;
        Amount?: {
            Value?: number | undefined;
            Exp?: number | undefined;
        } | undefined;
        Price?: number | undefined;
        Denom1?: {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } | undefined;
        Denom2?: {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } | undefined;
        Side?: Side | undefined;
        BlockTime?: Date | undefined;
        MetaData?: {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } | undefined;
        TXID?: string | undefined;
        BlockHeight?: number | undefined;
        Enriched?: boolean | undefined;
        USD?: number | undefined;
        Inverted?: boolean | undefined;
    } & {
        Account?: string | undefined;
        OrderID?: string | undefined;
        Sequence?: number | undefined;
        Amount?: ({
            Value?: number | undefined;
            Exp?: number | undefined;
        } & {
            Value?: number | undefined;
            Exp?: number | undefined;
        } & { [K_5 in Exclude<keyof I_1["Amount"], keyof Decimal>]: never; }) | undefined;
        Price?: number | undefined;
        Denom1?: ({
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
        } & { [K_6 in Exclude<keyof I_1["Denom1"], keyof Denom>]: never; }) | undefined;
        Denom2?: ({
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
        } & { [K_7 in Exclude<keyof I_1["Denom2"], keyof Denom>]: never; }) | undefined;
        Side?: Side | undefined;
        BlockTime?: Date | undefined;
        MetaData?: ({
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & { [K_8 in Exclude<keyof I_1["MetaData"], keyof MetaData>]: never; }) | undefined;
        TXID?: string | undefined;
        BlockHeight?: number | undefined;
        Enriched?: boolean | undefined;
        USD?: number | undefined;
        Inverted?: boolean | undefined;
    } & { [K_9 in Exclude<keyof I_1, keyof Trade>]: never; }>(object: I_1): Trade;
};
export declare const Trades: {
    encode(message: Trades, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Trades;
    fromJSON(object: any): Trades;
    toJSON(message: Trades): unknown;
    create<I extends {
        Trades?: {
            Account?: string | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            Amount?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Price?: number | undefined;
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Side?: Side | undefined;
            BlockTime?: Date | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
            USD?: number | undefined;
            Inverted?: boolean | undefined;
        }[] | undefined;
    } & {
        Trades?: ({
            Account?: string | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            Amount?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Price?: number | undefined;
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Side?: Side | undefined;
            BlockTime?: Date | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
            USD?: number | undefined;
            Inverted?: boolean | undefined;
        }[] & ({
            Account?: string | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            Amount?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Price?: number | undefined;
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Side?: Side | undefined;
            BlockTime?: Date | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
            USD?: number | undefined;
            Inverted?: boolean | undefined;
        } & {
            Account?: string | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            Amount?: ({
                Value?: number | undefined;
                Exp?: number | undefined;
            } & {
                Value?: number | undefined;
                Exp?: number | undefined;
            } & { [K in Exclude<keyof I["Trades"][number]["Amount"], keyof Decimal>]: never; }) | undefined;
            Price?: number | undefined;
            Denom1?: ({
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
            } & { [K_1 in Exclude<keyof I["Trades"][number]["Denom1"], keyof Denom>]: never; }) | undefined;
            Denom2?: ({
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
            } & { [K_2 in Exclude<keyof I["Trades"][number]["Denom2"], keyof Denom>]: never; }) | undefined;
            Side?: Side | undefined;
            BlockTime?: Date | undefined;
            MetaData?: ({
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & { [K_3 in Exclude<keyof I["Trades"][number]["MetaData"], keyof MetaData>]: never; }) | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
            USD?: number | undefined;
            Inverted?: boolean | undefined;
        } & { [K_4 in Exclude<keyof I["Trades"][number], keyof Trade>]: never; })[] & { [K_5 in Exclude<keyof I["Trades"], keyof {
            Account?: string | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            Amount?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Price?: number | undefined;
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Side?: Side | undefined;
            BlockTime?: Date | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
            USD?: number | undefined;
            Inverted?: boolean | undefined;
        }[]>]: never; }) | undefined;
    } & { [K_6 in Exclude<keyof I, "Trades">]: never; }>(base?: I | undefined): Trades;
    fromPartial<I_1 extends {
        Trades?: {
            Account?: string | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            Amount?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Price?: number | undefined;
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Side?: Side | undefined;
            BlockTime?: Date | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
            USD?: number | undefined;
            Inverted?: boolean | undefined;
        }[] | undefined;
    } & {
        Trades?: ({
            Account?: string | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            Amount?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Price?: number | undefined;
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Side?: Side | undefined;
            BlockTime?: Date | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
            USD?: number | undefined;
            Inverted?: boolean | undefined;
        }[] & ({
            Account?: string | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            Amount?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Price?: number | undefined;
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Side?: Side | undefined;
            BlockTime?: Date | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
            USD?: number | undefined;
            Inverted?: boolean | undefined;
        } & {
            Account?: string | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            Amount?: ({
                Value?: number | undefined;
                Exp?: number | undefined;
            } & {
                Value?: number | undefined;
                Exp?: number | undefined;
            } & { [K_7 in Exclude<keyof I_1["Trades"][number]["Amount"], keyof Decimal>]: never; }) | undefined;
            Price?: number | undefined;
            Denom1?: ({
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
            } & { [K_8 in Exclude<keyof I_1["Trades"][number]["Denom1"], keyof Denom>]: never; }) | undefined;
            Denom2?: ({
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
            } & { [K_9 in Exclude<keyof I_1["Trades"][number]["Denom2"], keyof Denom>]: never; }) | undefined;
            Side?: Side | undefined;
            BlockTime?: Date | undefined;
            MetaData?: ({
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & { [K_10 in Exclude<keyof I_1["Trades"][number]["MetaData"], keyof MetaData>]: never; }) | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
            USD?: number | undefined;
            Inverted?: boolean | undefined;
        } & { [K_11 in Exclude<keyof I_1["Trades"][number], keyof Trade>]: never; })[] & { [K_12 in Exclude<keyof I_1["Trades"], keyof {
            Account?: string | undefined;
            OrderID?: string | undefined;
            Sequence?: number | undefined;
            Amount?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            Price?: number | undefined;
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Side?: Side | undefined;
            BlockTime?: Date | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            TXID?: string | undefined;
            BlockHeight?: number | undefined;
            Enriched?: boolean | undefined;
            USD?: number | undefined;
            Inverted?: boolean | undefined;
        }[]>]: never; }) | undefined;
    } & { [K_13 in Exclude<keyof I_1, "Trades">]: never; }>(object: I_1): Trades;
};
export declare const TradePair: {
    encode(message: TradePair, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TradePair;
    fromJSON(object: any): TradePair;
    toJSON(message: TradePair): unknown;
    create<I extends {
        Denom1?: {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } | undefined;
        Denom2?: {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } | undefined;
        MetaData?: {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } | undefined;
        PriceTick?: {
            Value?: number | undefined;
            Exp?: number | undefined;
        } | undefined;
        QuantityStep?: number | undefined;
    } & {
        Denom1?: ({
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
        } & { [K in Exclude<keyof I["Denom1"], keyof Denom>]: never; }) | undefined;
        Denom2?: ({
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
        } & { [K_1 in Exclude<keyof I["Denom2"], keyof Denom>]: never; }) | undefined;
        MetaData?: ({
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & { [K_2 in Exclude<keyof I["MetaData"], keyof MetaData>]: never; }) | undefined;
        PriceTick?: ({
            Value?: number | undefined;
            Exp?: number | undefined;
        } & {
            Value?: number | undefined;
            Exp?: number | undefined;
        } & { [K_3 in Exclude<keyof I["PriceTick"], keyof Decimal>]: never; }) | undefined;
        QuantityStep?: number | undefined;
    } & { [K_4 in Exclude<keyof I, keyof TradePair>]: never; }>(base?: I | undefined): TradePair;
    fromPartial<I_1 extends {
        Denom1?: {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } | undefined;
        Denom2?: {
            Currency?: string | undefined;
            Issuer?: string | undefined;
            Precision?: number | undefined;
            IsIBC?: boolean | undefined;
            Denom?: string | undefined;
            Name?: string | undefined;
            Description?: string | undefined;
            Icon?: string | undefined;
        } | undefined;
        MetaData?: {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } | undefined;
        PriceTick?: {
            Value?: number | undefined;
            Exp?: number | undefined;
        } | undefined;
        QuantityStep?: number | undefined;
    } & {
        Denom1?: ({
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
        } & { [K_5 in Exclude<keyof I_1["Denom1"], keyof Denom>]: never; }) | undefined;
        Denom2?: ({
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
        } & { [K_6 in Exclude<keyof I_1["Denom2"], keyof Denom>]: never; }) | undefined;
        MetaData?: ({
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & { [K_7 in Exclude<keyof I_1["MetaData"], keyof MetaData>]: never; }) | undefined;
        PriceTick?: ({
            Value?: number | undefined;
            Exp?: number | undefined;
        } & {
            Value?: number | undefined;
            Exp?: number | undefined;
        } & { [K_8 in Exclude<keyof I_1["PriceTick"], keyof Decimal>]: never; }) | undefined;
        QuantityStep?: number | undefined;
    } & { [K_9 in Exclude<keyof I_1, keyof TradePair>]: never; }>(object: I_1): TradePair;
};
export declare const TradePairs: {
    encode(message: TradePairs, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TradePairs;
    fromJSON(object: any): TradePairs;
    toJSON(message: TradePairs): unknown;
    create<I extends {
        TradePairs?: {
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            PriceTick?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            QuantityStep?: number | undefined;
        }[] | undefined;
        Offset?: number | undefined;
    } & {
        TradePairs?: ({
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            PriceTick?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            QuantityStep?: number | undefined;
        }[] & ({
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            PriceTick?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            QuantityStep?: number | undefined;
        } & {
            Denom1?: ({
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
            } & { [K in Exclude<keyof I["TradePairs"][number]["Denom1"], keyof Denom>]: never; }) | undefined;
            Denom2?: ({
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
            } & { [K_1 in Exclude<keyof I["TradePairs"][number]["Denom2"], keyof Denom>]: never; }) | undefined;
            MetaData?: ({
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & { [K_2 in Exclude<keyof I["TradePairs"][number]["MetaData"], keyof MetaData>]: never; }) | undefined;
            PriceTick?: ({
                Value?: number | undefined;
                Exp?: number | undefined;
            } & {
                Value?: number | undefined;
                Exp?: number | undefined;
            } & { [K_3 in Exclude<keyof I["TradePairs"][number]["PriceTick"], keyof Decimal>]: never; }) | undefined;
            QuantityStep?: number | undefined;
        } & { [K_4 in Exclude<keyof I["TradePairs"][number], keyof TradePair>]: never; })[] & { [K_5 in Exclude<keyof I["TradePairs"], keyof {
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            PriceTick?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            QuantityStep?: number | undefined;
        }[]>]: never; }) | undefined;
        Offset?: number | undefined;
    } & { [K_6 in Exclude<keyof I, keyof TradePairs>]: never; }>(base?: I | undefined): TradePairs;
    fromPartial<I_1 extends {
        TradePairs?: {
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            PriceTick?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            QuantityStep?: number | undefined;
        }[] | undefined;
        Offset?: number | undefined;
    } & {
        TradePairs?: ({
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            PriceTick?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            QuantityStep?: number | undefined;
        }[] & ({
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            PriceTick?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            QuantityStep?: number | undefined;
        } & {
            Denom1?: ({
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
            } & { [K_7 in Exclude<keyof I_1["TradePairs"][number]["Denom1"], keyof Denom>]: never; }) | undefined;
            Denom2?: ({
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
            } & { [K_8 in Exclude<keyof I_1["TradePairs"][number]["Denom2"], keyof Denom>]: never; }) | undefined;
            MetaData?: ({
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & { [K_9 in Exclude<keyof I_1["TradePairs"][number]["MetaData"], keyof MetaData>]: never; }) | undefined;
            PriceTick?: ({
                Value?: number | undefined;
                Exp?: number | undefined;
            } & {
                Value?: number | undefined;
                Exp?: number | undefined;
            } & { [K_10 in Exclude<keyof I_1["TradePairs"][number]["PriceTick"], keyof Decimal>]: never; }) | undefined;
            QuantityStep?: number | undefined;
        } & { [K_11 in Exclude<keyof I_1["TradePairs"][number], keyof TradePair>]: never; })[] & { [K_12 in Exclude<keyof I_1["TradePairs"], keyof {
            Denom1?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            Denom2?: {
                Currency?: string | undefined;
                Issuer?: string | undefined;
                Precision?: number | undefined;
                IsIBC?: boolean | undefined;
                Denom?: string | undefined;
                Name?: string | undefined;
                Description?: string | undefined;
                Icon?: string | undefined;
            } | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            PriceTick?: {
                Value?: number | undefined;
                Exp?: number | undefined;
            } | undefined;
            QuantityStep?: number | undefined;
        }[]>]: never; }) | undefined;
        Offset?: number | undefined;
    } & { [K_13 in Exclude<keyof I_1, keyof TradePairs>]: never; }>(object: I_1): TradePairs;
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
