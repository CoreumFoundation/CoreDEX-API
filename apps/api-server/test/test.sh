#!/bin/bash

# Global variable to store the current host
HOST="http://localhost:8080/api"
WS_HOST="ws://localhost:8080/api"

# used for TX submit
senderMnemonic="silk loop drastic novel taste project mind dragon shock outside stove patrol immense car collect winter melody pizza all deputy kid during style ribbon"


function setNetworkDevnet() {
    NETWORK=devnet
    SYMBOL=nor-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43_alb-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43
}

function setNetworkTestnet() {
    NETWORK=testnet
    SYMBOL=nor-testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57_alb-testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57
}

# Function to set the host to localhost
set_host_local() {
    HOST="http://localhost:8080/api"
    WS_HOST="ws://localhost:8080/api"
    echo "Host set to localhost"
}

# Function to set the host to https://coredex.test.coreum.dev
set_host_coreum() {
    HOST="https://coredex.test.coreum.dev/api"
    WS_HOST="wss://coredex.test.coreum.dev/api"
    echo "Host set to coredex.test.coreum.dev"
}

# Function for GET /ohlc
get_ohlc() {
    local symbol=${SYMBOL}
    local periods=("1m" "3m" "5m" "15m" "30m" "1h" "3h" "6h" "12h" "1d" "3d" "1w")

    # Calculate 'to' as the current time and 'from' as 24 hours before 'to'
    local to=$(date +%s)
    local from=$((to - 86400))

    # URL encode the symbol
    local encoded_symbol=$(echo -n "$symbol" | jq -sRr @uri)

    # Iterate over periods and make the API call for each
    for period in "${periods[@]}"; do
        echo "Calling API for OHLC with period ${period}"
        curl "$HOST/ohlc?symbol=${encoded_symbol}&period=${period}&from=${from}&to=${to}" \
            --header "Network: ${NETWORK}"
        echo -e "\n"
    done
}

get_trades_with_account() {
    local symbol=${SYMBOL}
    local account="devcore1dj9yphkprdsuk6s4mgnfhnq5c39zf499nknkna"
    local to=1734462880
    local from=$((to - 86400))
    local encoded_symbol=$(echo -n "$symbol" | jq -sRr @uri)

    echo "Calling API for trades with account ${account}"
    curl "${HOST}/trades?symbol=${encoded_symbol}&from=${from}&to=${to}&account=${account}" \
        --header "Network: devnet"
    echo -e "\n"
}

get_trades_without_account() {
    local symbol=${SYMBOL}
    local to=$(date +%s)
    local from=$((to - 86400))
    local encoded_symbol=$(echo -n "$symbol" | jq -sRr @uri)

    echo "Calling API for trades without account"
    curl "${HOST}/trades?symbol=${encoded_symbol}&from=${from}&to=${to}&side=1" \
        --header "Network: ${NETWORK}"
    echo -e "\n"
}

get_trades_without_account_inverted() {
    local symbol="alb-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43_nor-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43"
    local to=$(date +%s)
    to=1740700197
    local from=$((to - 86400))
    local encoded_symbol=$(echo -n "$symbol" | jq -sRr @uri)

    echo "Calling API for trades without account currencies inverted"
    curl "${HOST}/trades?symbol=${encoded_symbol}&from=${from}&to=${to}&side=1" \
        --header "Network: ${NETWORK}"
    echo -e "\n"
}

# Function for GET /tickers
get_tickers() {
    local symbols=(${SYMBOL})
    local json_symbols=$(printf '%s\n' "${symbols[@]}" | jq -R . | jq -s .)
    local encoded_symbols=""
    if base64 --help 2>&1 | grep -q -- "-w"; then
        # GNU base64 (Linux)
        encoded_symbols=$(echo -n "$json_symbols" | base64 -w 0)
    else
        # BSD base64 (macOS)
        encoded_symbols=$(echo -n "$json_symbols" | base64 | tr -d '\n')
    fi
    echo "Calling API for tickers"
    curl -H "Network: ${NETWORK}" \
         -X "GET" "${HOST}/tickers?symbols=${encoded_symbols}"
    echo -e "\n"
}

# Function for GET /currencies
get_currencies() {
    echo "Calling GET /currencies"
    curl -H "Network: ${NETWORK}" \
         -X "GET" "${HOST}/currencies"
}

# Function for GET /ws
# hint: brew install websocat
get_ws() {
    echo "Calling GET /ws"
   go run ./ws/main.go
}

# Function for POST /order/create
post_order_create() {
    echo "Calling POST /order/create"
    post_order_create_dry
}

post_order_create_dry() {
    curl -s -H "Network: ${NETWORK}" \
    -X "POST" "${HOST}/order/create" \
    -d '{
    "Sender": "devcore1fksu90amj2qgydf43dm2qf6m2dl4szjtx6j5q8",
    "Type": 1,
    "BaseDenom": "alb-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
    "QuoteDenom": "nor-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
    "Price": "0.25",
    "Quantity": "1000",
    "Side": 1,
    "GoodTil": {
        "GoodTilBlockTime": "2025-12-30T12:00:00Z"
    },
    "TimeInForce": 1
    }'
}

# Function for POST /order/submit
post_order_submit() {
    echo "Calling POST /order/submit: Creates an order using order/create, signs with the app in test/sign"
    # Add your curl or wget command here
    response=$(post_order_create_dry)
    tx_bytes=$(echo "${response}" | jq -r '.TXBytes')
    signedTX=$(go run sign/main.go ${tx_bytes} ${senderMnemonic})
    # Split by space in signed base64 tx and the account (present in output for verification purposes):
    IFS=' ' read -r -a txSigned <<< "$signedTX"

    echo "Submitting the signed tx to the order/submit endpoint: ${txSigned[0]} for account ${txSigned[1]}"

    # Submit the txSigned[0] to the order/submit endpoint:
    curl -H "Network: ${NETWORK}" -X "POST" \
    "${HOST}/order/submit" \
    -d '{"TX": "'${txSigned[0]}'"}'
}

# Function for GET /order/orderbook
get_order_orderbook() {
    echo "Calling GET /order/orderbook"
    # Add your curl or wget command here
    curl -H "Network: ${NETWORK}" \
    -X "GET" "${HOST}/order/orderbook?symbol=nor-testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57_alb-testcore1eyhq55grezrggrxs9eweml7nw7alkd8hv9vt57"
}

# Function for GET /order/orderbook
get_order_orderbook_for_account() {
    echo "Calling GET /order/orderbook"
    # Add your curl or wget command here
    curl -H "Network: ${NETWORK}" \
    -X "GET" "${HOST}/order/orderbook?symbol=nor-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43_alb-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43&account=devcore1fksu90amj2qgydf43dm2qf6m2dl4szjtx6j5q8"
}

get_wallet_assets() {
    echo "Calling GET /wallet/assets"
    curl -H "Network: ${NETWORK}" \
    -X "GET" "${HOST}/wallet/assets?address=devcore1fksu90amj2qgydf43dm2qf6m2dl4szjtx6j5q8"
}

post_order_cancel() {
    echo "Calling POST /order/cancel"
    # Add your curl or wget command here
    response=$(curl -H "Network: ${NETWORK}" \
    -X "POST" "${HOST}/order/cancel" \
    -d '{
    "Sender": "devcore1fksu90amj2qgydf43dm2qf6m2dl4szjtx6j5q8",
    "OrderID": "d92cc4a6-24f0-42e6-bb2f-6a69bef1f2ce"
    }')
    echo "Cancel response: ${response}"
    # Sign the TXBytes with the app in test/sign
    # tx_bytes=$(echo "${response}" | jq -r '.TXBytes')
    # signedTX=$(go run sign/main.go ${tx_bytes} ${senderMnemonic})
    # # Split by space in signed base64 tx and the account (present in output for verification purposes):
    # IFS=' ' read -r -a txSigned <<< "$signedTX"

    # echo "Submitting the signed cancel tx to the order/submit endpoint: ${txSigned[0]} for account ${txSigned[1]}"

    # # Submit the txSigned[0] to the order/submit endpoint:
    # response=$(curl -H "Network: devnet" -X "POST" \
    # "${HOST}/order/submit" \
    # -d '{"TX": "'${txSigned[0]}'"}')
    # echo "Cancel result: ${response}"
}

get_market() {
    echo "Calling GET /market"
    curl -H "Network: ${NETWORK}" \
    -X "GET" "${HOST}/market?symbol=alb-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43_nor-devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43"
}

# Display menu
show_menu() {
    printf "│ %-30s │ %-30s │\n" "a. Set host to localhost" "7. POST /order/submit"
    printf "│ %-30s │ %-30s │\n" "b. Set host to api.coreum.com" "8. POST /order/cancel"
    printf "│ %-30s │ %-30s │\n" "c. Set network to devnet" "9. GET /order/orderbook"
    printf "│ %-30s │ %-30s │\n" "d. Set network to testnet" "10. GET /orderbook for account"
    printf "│ %-30s │ %-30s │\n" "0. Exit" "11. GET /ws"
    printf "│ %-30s │ %-30s │\n" "1. GET /ohlc" "12. GET /currencies"
    printf "│ %-30s │ %-30s │\n" "2. GET /trades with account" "13. GET /trades (inverted)"
    printf "│ %-30s │ %-30s │\n" "3. GET /trades without account" "14. GET /market"
    printf "│ %-30s │ %-30s │\n" "4. GET /tickers" ""
    printf "│ %-30s │ %-30s │\n" "5. GET /wallet/assets" ""
    printf "│ %-30s │ %-30s │\n" "6. POST /order/create" ""
}

# Main loop
while true; do
    show_menu
    read -p "Enter your choice: " choice
    case $choice in
        a) set_host_local ;;
        b) set_host_coreum ;;
        c) setNetworkDevnet ;;
        d) setNetworkTestnet ;;
        0) echo "Exiting..."; exit 0 ;;
        1) get_ohlc ;;
        2) get_trades_with_account ;;
        3) get_trades_without_account ;;
        4) get_tickers ;;
        5) get_wallet_assets ;;
        6) post_order_create ;;
        7) post_order_submit ;;
        8) post_order_cancel ;;
        9) get_order_orderbook ;;
        10) get_order_orderbook_for_account ;;
        11) get_ws ;;
        12) get_currencies ;;
        13) get_trades_without_account_inverted ;;
        14) get_market ;;
        *) echo "Invalid choice, please try again." ;;
    esac
done