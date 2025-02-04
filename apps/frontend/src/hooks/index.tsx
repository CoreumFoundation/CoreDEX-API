import { Dispatch, SetStateAction, useEffect, useState } from "react";
import {
  autoUpdate,
  shift,
  Placement,
  computePosition,
} from "@floating-ui/react";
import { IChartingLibraryWidget } from "@/vendor/tradingview/charting_library/charting_library";
import { useStore } from "@/state/store";
import { getOverrides } from "@/components/TradingView/tools/config";

export enum TooltipPosition {
  TOP = "top",
  TOP_RIGHT = "top-end",
  TOP_LEFT = "top-start",
  RIGHT = "right",
  BOTTOM = "bottom",
  BOTTOM_LEFT = "bottom-start",
  BOTTOM_RIGHT = "bottom-end",
  LEFT = "left",
}

export const useTooltip = () => {
  const [isShown, setIsShown] = useState(false);

  const showTooltip = (
    element: HTMLElement,
    content: any,
    position: TooltipPosition,
    width?: string
  ) => {
    let tooltip = document.getElementById("tooltip-message");

    if (!tooltip) {
      tooltip = document.createElement("div");
      tooltip.id = "tooltip-message";
      tooltip.style.width = width ? width : "fit-content";
    }
    tooltip.innerHTML = content;

    function wrapper(referenceEl: HTMLElement, floatingEl: HTMLElement) {
      return computePosition(referenceEl, floatingEl, {
        placement: position as Placement,
        middleware: [shift()],
      });
    }
    const cleanup = autoUpdate(element, tooltip, () => {
      if (!tooltip) return;
      wrapper(element, tooltip).then(({ x, y }) => {
        Object.assign(tooltip!.style, {
          left: `${x}px`,
          top: `${y}px`,
        });
      });
    });
    document.getElementById("root")!.appendChild(tooltip!);
    cleanup();
  };

  const hideTooltip = () => {
    const tooltip = document.getElementById("tooltip-message");
    if (tooltip) {
      tooltip.style.opacity = "0"; // Fade out
      // Delay removal to allow for transition
      setIsShown(false);
      tooltip.remove();
    }
  };

  return {
    showTooltip,
    hideTooltip,
    isShown,
  };
};

// tradingview
export function useMountChart(mountChart: () => void) {
  const { market, tickers } = useStore();
  const [chartReady, setReady] = useState(false);

  useEffect(() => {
    setReady(false);

    if (tickers) {
      setTimeout(() => {
        mountChart();
      }, 300);
    }
  }, [market]);

  return { chartReady, setReady };
}

export function useSaveAndClear(
  mountChart: () => void,
  setReady: Dispatch<SetStateAction<boolean>>
) {
  const { market } = useStore();

  const { pushNotification } = useStore();
  const [clearable, setClearable] = useState(
    market?.pair_symbol && localStorage.getItem(market.pair_symbol)
      ? true
      : false
  );

  useEffect(() => {
    if (!clearable) mountChart();
  }, [clearable]);

  return {
    clearable,
    saveChart: () => {
      window.tvWidget.save((chartData: any) => {
        if (market?.pair_symbol)
          window.localStorage.setItem(
            market.pair_symbol,
            JSON.stringify(chartData)
          );
      });
      pushNotification({
        message: "_chartSaved",
        type: "success",
      });
      setClearable(true);
    },
    clearChart: () => {
      if (market?.pair_symbol) {
        window.localStorage.removeItem(market.pair_symbol);
        setClearable(false);
        setReady(false);
      }
    },
  };
}

export function useChartTheme(
  chartReady: boolean | Dispatch<SetStateAction<boolean>>
) {
  useEffect(() => {
    if (chartReady) {
      const widget = window.tvWidget as IChartingLibraryWidget;
      widget
        .changeTheme("Dark")
        .then(() => {
          widget.applyOverrides(getOverrides("dark"));
        })
        .catch((e) => {
          console.error(e);
        });
    }
  }, [chartReady]);
}
