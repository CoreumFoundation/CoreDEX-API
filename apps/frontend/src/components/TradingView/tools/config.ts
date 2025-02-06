import themes from "./theme";

const time_frames = [
  {
    text: "1y",
    resolution: "1W",
  },
  {
    text: "1m",
    resolution: "1D",
  },
  {
    text: "1d",
    resolution: "30",
  },
  {
    text: "12h",
    resolution: "15",
  },
  {
    text: "1h",
    resolution: "1",
  },
];

const drawings_access = {
  type: "black",
  tools: [{ name: "Regression Trend" }],
};

const studies_overrides = {
  "moving average exponential.plot.color": "#efff00",
  "moving average exponential.plot.linewidth": 3,
  "volume.volume.color.0": "rgba(216, 29, 60, 0.25)", // red volume bars
  "volume.volume.color.1": "rgba(37, 214, 149, 0.25)", // green volume bars
};

const disabled_features = [
  "header_symbol_search",
  "header_interval_dialog_button",
  "header_settings",
  "header_compare",
  "header_undo_redo",
  "header_indicators",
  "header_screenshot",
  "header_widget",
  "compare_symbol",
  "context_menus",
  "volume_force_overlay",
  "use_localstorage_for_settings",
  "study_templates",
  "display_market_status",
  "header_saveload",
];

const enabled_features: string[] = [
  // "move_logo_to_main_pane",
  // "header_chart_type",
  // "header_resolutions",
  // "header_fullscreen_button",
];

const misc = {
  width: "100%",
  library_path: window.location.href.split("/")[0] + "/trade/charting_library/",
  container: "chartContainer",
  allow_symbol_change: false,
  load_last_chart: false,
};

export const DEFAULT_CONFIGS = {
  time_frames,
  drawings_access,
  studies_overrides,
  disabled_features,
  enabled_features,
  toolbar_bg: "#101216",
  layout: {
    attributionLogo: false,
  },
  ...misc,
};

export function getOverrides(theme: string) {
  return {
    "paneProperties.backgroundType": "solid",
    "paneProperties.background": themes[theme].colors.chart.background,
    "paneProperties.vertGridProperties.color":
      themes[theme].colors.chart.grid_lines,
    "paneProperties.horzGridProperties.color":
      themes[theme].colors.chart.grid_lines,
    "scalesProperties.textColor": themes[theme].colors.chart.scale_lines,
    "mainSeriesProperties.candleStyle.upColor":
      themes[theme].colors.chart.up_lines,
    "mainSeriesProperties.candleStyle.downColor":
      themes[theme].colors.chart.down_lines,
    "mainSeriesProperties.candleStyle.borderUpColor":
      themes[theme].colors.chart.up_lines,
    "mainSeriesProperties.candleStyle.borderDownColor":
      themes[theme].colors.chart.down_lines,
  };
}
