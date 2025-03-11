export declare const protobufPackage = "orderproperties";
/** Side is order side. */
export declare enum Side {
    /** SIDE_UNSPECIFIED - SIDE_UNSPECIFIED reserves the default value, to protect against unexpected settings. */
    SIDE_UNSPECIFIED = 0,
    /** SIDE_BUY - SIDE_BUY means that the order is to buy base_denom quantity with the price. */
    SIDE_BUY = 1,
    /** SIDE_SELL - SIDE_SELL means that the order is to sell base_denom quantity with the price. */
    SIDE_SELL = 2,
    UNRECOGNIZED = -1
}
export declare function sideFromJSON(object: any): Side;
export declare function sideToJSON(object: Side): string;
