import _m0 from "protobufjs/minimal";
export declare const protobufPackage = "denom";
export interface Denom {
    Currency: string;
    Issuer: string;
    Precision?: number | undefined;
    IsIBC: boolean;
    /** Could be handy for IBC */
    Denom: string;
    /** Additional fields, make it possible to use the denom as the currency storage for display purposes */
    Name?: string | undefined;
    Description?: string | undefined;
    Icon?: string | undefined;
}
export declare const Denom: {
    encode(message: Denom, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Denom;
    fromJSON(object: any): Denom;
    toJSON(message: Denom): unknown;
    create<I extends {
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
    } & { [K in Exclude<keyof I, keyof Denom>]: never; }>(base?: I | undefined): Denom;
    fromPartial<I_1 extends {
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
    } & { [K_1 in Exclude<keyof I_1, keyof Denom>]: never; }>(object: I_1): Denom;
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
