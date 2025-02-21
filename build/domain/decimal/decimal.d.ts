import _m0 from "protobufjs/minimal";
export declare const protobufPackage = "decimal";
export interface Decimal {
    Value: number;
    Exp: number;
}
export declare const Decimal: {
    encode(message: Decimal, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Decimal;
    fromJSON(object: any): Decimal;
    toJSON(message: Decimal): unknown;
    create<I extends {
        Value?: number | undefined;
        Exp?: number | undefined;
    } & {
        Value?: number | undefined;
        Exp?: number | undefined;
    } & { [K in Exclude<keyof I, keyof Decimal>]: never; }>(base?: I | undefined): Decimal;
    fromPartial<I_1 extends {
        Value?: number | undefined;
        Exp?: number | undefined;
    } & {
        Value?: number | undefined;
        Exp?: number | undefined;
    } & { [K_1 in Exclude<keyof I_1, keyof Decimal>]: never; }>(object: I_1): Decimal;
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
