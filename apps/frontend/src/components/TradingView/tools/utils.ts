import { getOHLC } from "@/services/api";
import { ChartFeedBarsParams } from "@/types/market";

export const fetchOLHCData = async (
  symbol: string,
  resolution: string,
  from: number,
  to: number
) => {
  const params: ChartFeedBarsParams = {
    symbol,
    from: from.toFixed(0),
    to: to.toFixed(0),
    period: resolveResolution(resolution),
  };
  const response = await getOHLC(symbol, params.period, params.from, params.to);

  if (response.status === 200 && response.data.length > 0) {
    const feedData = response.data;
    return feedData.map((el) => ({
      time: Number(el[0] * 1000),
      open: Number(el[1]),
      high: Number(el[2]),
      low: Number(el[3]),
      close: Number(el[4]),
      volume: Number(el[5]),
    }));
  } else return [];
};

export function resolveResolution(resolution: string) {
  switch (resolution) {
    case "1":
      return "1m";
    case "3":
      return "3m";
    case "5":
      return "5m";
    case "15":
      return "15m";
    case "30":
      return "30m";
    case "1h":
      return "1h";
    case "3h":
      return "3h";
    case "6h":
      return "6h";
    case "12h":
      return "12h";
    case "1D":
      return "1d";
    case "3D":
      return "3d";
    case "1W":
      return "1w";
    default:
      return "1h";
  }
}

export const SUPPORTED_RESOLUTIONS = [
  "1",
  "3",
  "5",
  "15",
  "30",
  "60",
  "120",
  "180",
  "360",
  "720",
  "1D",
  "3D",
  "1W",
];
