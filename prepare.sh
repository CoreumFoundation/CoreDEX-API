#!/usr/bin/env sh

set -xe

for variable in \
    VITE_ENV_BASE_API \
    VITE_ENV_BASE_WS \
    VITE_ENV_MODE \
    VITE_ENV_WS \
    VITE_ALLOWED_HOST ; do
  eval value="\$${variable}"
  echo ${value}
  sed -E -i "s|\{\{$variable\}\}|$value|g" /usr/share/nginx/html/index.html;
done

nginx -g "daemon off;"
