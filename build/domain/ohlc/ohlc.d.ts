import _m0 from "protobufjs/minimal";
import { MetaData } from "../metadata/metadata";
export declare const protobufPackage = "ohlc";
export declare enum PeriodType {
    PERIOD_TYPE_DO_NOT_USE = 0,
    PERIOD_TYPE_MINUTE = 1,
    PERIOD_TYPE_HOUR = 2,
    PERIOD_TYPE_DAY = 3,
    PERIOD_TYPE_WEEK = 4,
    UNRECOGNIZED = -1
}
export declare function periodTypeFromJSON(object: any): PeriodType;
export declare function periodTypeToJSON(object: PeriodType): string;
export interface OHLCs {
    OHLCs: OHLC[];
}
export interface OHLC {
    Symbol: string;
    Timestamp: Date | undefined;
    Open: number;
    High: number;
    Low: number;
    Close: number;
    Volume: number;
    NumberOfTrades: number;
    Period: Period | undefined;
    USDValue?: number | undefined;
    QuoteVolume: number;
    MetaData: MetaData | undefined;
    /** When was the open time record created: Used for out of order trade processing */
    OpenTime: Date | undefined;
    /** When was the close time record created: Used for out of order trade processing */
    CloseTime: Date | undefined;
}
export interface Period {
    PeriodType: PeriodType;
    /** The duration of the indicated period (e.g 1 minute, 3 minutes, etc) */
    Duration: number;
}
export declare const OHLCs: {
    encode(message: OHLCs, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): OHLCs;
    fromJSON(object: any): OHLCs;
    toJSON(message: OHLCs): unknown;
    create<I extends {
        OHLCs?: {
            Symbol?: string | undefined;
            Timestamp?: Date | undefined;
            Open?: number | undefined;
            High?: number | undefined;
            Low?: number | undefined;
            Close?: number | undefined;
            Volume?: number | undefined;
            NumberOfTrades?: number | undefined;
            Period?: {
                PeriodType?: PeriodType | undefined;
                Duration?: number | undefined;
            } | undefined;
            USDValue?: number | undefined;
            QuoteVolume?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            OpenTime?: Date | undefined;
            CloseTime?: Date | undefined;
        }[] | undefined;
    } & {
        OHLCs?: ({
            Symbol?: string | undefined;
            Timestamp?: Date | undefined;
            Open?: number | undefined;
            High?: number | undefined;
            Low?: number | undefined;
            Close?: number | undefined;
            Volume?: number | undefined;
            NumberOfTrades?: number | undefined;
            Period?: {
                PeriodType?: PeriodType | undefined;
                Duration?: number | undefined;
            } | undefined;
            USDValue?: number | undefined;
            QuoteVolume?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            OpenTime?: Date | undefined;
            CloseTime?: Date | undefined;
        }[] & ({
            Symbol?: string | undefined;
            Timestamp?: Date | undefined;
            Open?: number | undefined;
            High?: number | undefined;
            Low?: number | undefined;
            Close?: number | undefined;
            Volume?: number | undefined;
            NumberOfTrades?: number | undefined;
            Period?: {
                PeriodType?: PeriodType | undefined;
                Duration?: number | undefined;
            } | undefined;
            USDValue?: number | undefined;
            QuoteVolume?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            OpenTime?: Date | undefined;
            CloseTime?: Date | undefined;
        } & {
            Symbol?: string | undefined;
            Timestamp?: Date | undefined;
            Open?: number | undefined;
            High?: number | undefined;
            Low?: number | undefined;
            Close?: number | undefined;
            Volume?: number | undefined;
            NumberOfTrades?: number | undefined;
            Period?: ({
                PeriodType?: PeriodType | undefined;
                Duration?: number | undefined;
            } & {
                PeriodType?: PeriodType | undefined;
                Duration?: number | undefined;
            } & { [K in Exclude<keyof I["OHLCs"][number]["Period"], keyof Period>]: never; }) | undefined;
            USDValue?: number | undefined;
            QuoteVolume?: number | undefined;
            MetaData?: ({
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & { [K_1 in Exclude<keyof I["OHLCs"][number]["MetaData"], keyof MetaData>]: never; }) | undefined;
            OpenTime?: Date | undefined;
            CloseTime?: Date | undefined;
        } & { [K_2 in Exclude<keyof I["OHLCs"][number], keyof OHLC>]: never; })[] & { [K_3 in Exclude<keyof I["OHLCs"], keyof {
            Symbol?: string | undefined;
            Timestamp?: Date | undefined;
            Open?: number | undefined;
            High?: number | undefined;
            Low?: number | undefined;
            Close?: number | undefined;
            Volume?: number | undefined;
            NumberOfTrades?: number | undefined;
            Period?: {
                PeriodType?: PeriodType | undefined;
                Duration?: number | undefined;
            } | undefined;
            USDValue?: number | undefined;
            QuoteVolume?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            OpenTime?: Date | undefined;
            CloseTime?: Date | undefined;
        }[]>]: never; }) | undefined;
    } & { [K_4 in Exclude<keyof I, "OHLCs">]: never; }>(base?: I | undefined): OHLCs;
    fromPartial<I_1 extends {
        OHLCs?: {
            Symbol?: string | undefined;
            Timestamp?: Date | undefined;
            Open?: number | undefined;
            High?: number | undefined;
            Low?: number | undefined;
            Close?: number | undefined;
            Volume?: number | undefined;
            NumberOfTrades?: number | undefined;
            Period?: {
                PeriodType?: PeriodType | undefined;
                Duration?: number | undefined;
            } | undefined;
            USDValue?: number | undefined;
            QuoteVolume?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            OpenTime?: Date | undefined;
            CloseTime?: Date | undefined;
        }[] | undefined;
    } & {
        OHLCs?: ({
            Symbol?: string | undefined;
            Timestamp?: Date | undefined;
            Open?: number | undefined;
            High?: number | undefined;
            Low?: number | undefined;
            Close?: number | undefined;
            Volume?: number | undefined;
            NumberOfTrades?: number | undefined;
            Period?: {
                PeriodType?: PeriodType | undefined;
                Duration?: number | undefined;
            } | undefined;
            USDValue?: number | undefined;
            QuoteVolume?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            OpenTime?: Date | undefined;
            CloseTime?: Date | undefined;
        }[] & ({
            Symbol?: string | undefined;
            Timestamp?: Date | undefined;
            Open?: number | undefined;
            High?: number | undefined;
            Low?: number | undefined;
            Close?: number | undefined;
            Volume?: number | undefined;
            NumberOfTrades?: number | undefined;
            Period?: {
                PeriodType?: PeriodType | undefined;
                Duration?: number | undefined;
            } | undefined;
            USDValue?: number | undefined;
            QuoteVolume?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            OpenTime?: Date | undefined;
            CloseTime?: Date | undefined;
        } & {
            Symbol?: string | undefined;
            Timestamp?: Date | undefined;
            Open?: number | undefined;
            High?: number | undefined;
            Low?: number | undefined;
            Close?: number | undefined;
            Volume?: number | undefined;
            NumberOfTrades?: number | undefined;
            Period?: ({
                PeriodType?: PeriodType | undefined;
                Duration?: number | undefined;
            } & {
                PeriodType?: PeriodType | undefined;
                Duration?: number | undefined;
            } & { [K_5 in Exclude<keyof I_1["OHLCs"][number]["Period"], keyof Period>]: never; }) | undefined;
            USDValue?: number | undefined;
            QuoteVolume?: number | undefined;
            MetaData?: ({
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } & { [K_6 in Exclude<keyof I_1["OHLCs"][number]["MetaData"], keyof MetaData>]: never; }) | undefined;
            OpenTime?: Date | undefined;
            CloseTime?: Date | undefined;
        } & { [K_7 in Exclude<keyof I_1["OHLCs"][number], keyof OHLC>]: never; })[] & { [K_8 in Exclude<keyof I_1["OHLCs"], keyof {
            Symbol?: string | undefined;
            Timestamp?: Date | undefined;
            Open?: number | undefined;
            High?: number | undefined;
            Low?: number | undefined;
            Close?: number | undefined;
            Volume?: number | undefined;
            NumberOfTrades?: number | undefined;
            Period?: {
                PeriodType?: PeriodType | undefined;
                Duration?: number | undefined;
            } | undefined;
            USDValue?: number | undefined;
            QuoteVolume?: number | undefined;
            MetaData?: {
                Network?: import("../metadata/metadata").Network | undefined;
                UpdatedAt?: Date | undefined;
                CreatedAt?: Date | undefined;
            } | undefined;
            OpenTime?: Date | undefined;
            CloseTime?: Date | undefined;
        }[]>]: never; }) | undefined;
    } & { [K_9 in Exclude<keyof I_1, "OHLCs">]: never; }>(object: I_1): OHLCs;
};
export declare const OHLC: {
    encode(message: OHLC, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): OHLC;
    fromJSON(object: any): OHLC;
    toJSON(message: OHLC): unknown;
    create<I extends {
        Symbol?: string | undefined;
        Timestamp?: Date | undefined;
        Open?: number | undefined;
        High?: number | undefined;
        Low?: number | undefined;
        Close?: number | undefined;
        Volume?: number | undefined;
        NumberOfTrades?: number | undefined;
        Period?: {
            PeriodType?: PeriodType | undefined;
            Duration?: number | undefined;
        } | undefined;
        USDValue?: number | undefined;
        QuoteVolume?: number | undefined;
        MetaData?: {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } | undefined;
        OpenTime?: Date | undefined;
        CloseTime?: Date | undefined;
    } & {
        Symbol?: string | undefined;
        Timestamp?: Date | undefined;
        Open?: number | undefined;
        High?: number | undefined;
        Low?: number | undefined;
        Close?: number | undefined;
        Volume?: number | undefined;
        NumberOfTrades?: number | undefined;
        Period?: ({
            PeriodType?: PeriodType | undefined;
            Duration?: number | undefined;
        } & {
            PeriodType?: PeriodType | undefined;
            Duration?: number | undefined;
        } & { [K in Exclude<keyof I["Period"], keyof Period>]: never; }) | undefined;
        USDValue?: number | undefined;
        QuoteVolume?: number | undefined;
        MetaData?: ({
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & { [K_1 in Exclude<keyof I["MetaData"], keyof MetaData>]: never; }) | undefined;
        OpenTime?: Date | undefined;
        CloseTime?: Date | undefined;
    } & { [K_2 in Exclude<keyof I, keyof OHLC>]: never; }>(base?: I | undefined): OHLC;
    fromPartial<I_1 extends {
        Symbol?: string | undefined;
        Timestamp?: Date | undefined;
        Open?: number | undefined;
        High?: number | undefined;
        Low?: number | undefined;
        Close?: number | undefined;
        Volume?: number | undefined;
        NumberOfTrades?: number | undefined;
        Period?: {
            PeriodType?: PeriodType | undefined;
            Duration?: number | undefined;
        } | undefined;
        USDValue?: number | undefined;
        QuoteVolume?: number | undefined;
        MetaData?: {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } | undefined;
        OpenTime?: Date | undefined;
        CloseTime?: Date | undefined;
    } & {
        Symbol?: string | undefined;
        Timestamp?: Date | undefined;
        Open?: number | undefined;
        High?: number | undefined;
        Low?: number | undefined;
        Close?: number | undefined;
        Volume?: number | undefined;
        NumberOfTrades?: number | undefined;
        Period?: ({
            PeriodType?: PeriodType | undefined;
            Duration?: number | undefined;
        } & {
            PeriodType?: PeriodType | undefined;
            Duration?: number | undefined;
        } & { [K_3 in Exclude<keyof I_1["Period"], keyof Period>]: never; }) | undefined;
        USDValue?: number | undefined;
        QuoteVolume?: number | undefined;
        MetaData?: ({
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & {
            Network?: import("../metadata/metadata").Network | undefined;
            UpdatedAt?: Date | undefined;
            CreatedAt?: Date | undefined;
        } & { [K_4 in Exclude<keyof I_1["MetaData"], keyof MetaData>]: never; }) | undefined;
        OpenTime?: Date | undefined;
        CloseTime?: Date | undefined;
    } & { [K_5 in Exclude<keyof I_1, keyof OHLC>]: never; }>(object: I_1): OHLC;
};
export declare const Period: {
    encode(message: Period, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Period;
    fromJSON(object: any): Period;
    toJSON(message: Period): unknown;
    create<I extends {
        PeriodType?: PeriodType | undefined;
        Duration?: number | undefined;
    } & {
        PeriodType?: PeriodType | undefined;
        Duration?: number | undefined;
    } & { [K in Exclude<keyof I, keyof Period>]: never; }>(base?: I | undefined): Period;
    fromPartial<I_1 extends {
        PeriodType?: PeriodType | undefined;
        Duration?: number | undefined;
    } & {
        PeriodType?: PeriodType | undefined;
        Duration?: number | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof Period>]: never; }>(object: I_1): Period;
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
