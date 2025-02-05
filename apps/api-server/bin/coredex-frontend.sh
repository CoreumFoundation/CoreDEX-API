#!/bin/bash
helm upgrade --install "coredex-frontend" charts/frontend \
    --namespace $COREUM_NETWORK \
    --values $COREUM_NETWORK/backend/config/coredex-frontend.yaml