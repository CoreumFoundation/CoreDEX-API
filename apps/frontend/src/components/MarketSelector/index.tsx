import { useEffect, useRef, useState } from "react";
import Button from "../Button";
import Modal from "../Modal";
import Dropdown, { DropdownVariant } from "../Dropdown";
import { useStore } from "@/state/store";
import { Token } from "@/types/market";
import { getCurrencies } from "@/services/general";
import "./market-selector.scss";

const MarketSelector = () => {
  const { setCurrencies, currencies, market, setMarket } = useStore();
  const [isOpen, setIsOpen] = useState(false);

  const [baseToken, setBaseToken] = useState<Token | null>(
    market?.base ? market?.base : null
  );
  const [quoteToken, setQuoteToken] = useState<Token | null>(
    market?.counter ? market?.counter : null
  );

  const [openCreatePairModal, setOpenCreatePairModal] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchCurrencies = async () => {
      try {
        const data = await getCurrencies();
        if (data) {
          const currs = data.data.Currencies;
          setCurrencies(currs);
        }
      } catch (e) {
        console.log("ERROR GETTING CURRENCIES DATA >>", e);
        setCurrencies(null);
      }
    };
    fetchCurrencies();
  }, []);

  useEffect(() => {
    const handleOutsideClick = (event: MouseEvent) => {
      if (ref.current && !ref.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener("mousedown", handleOutsideClick);

    return () => {
      document.removeEventListener("mousedown", handleOutsideClick);
    };
  }, []);

  return (
    <div className="market-selector-container" ref={ref}>
      {/* Market Label Button */}
      <div
        className="market-label"
        onClick={() => {
          setOpenCreatePairModal(true);
        }}
      >
        <div className="market-label-selected">
          {market?.base.Denom.Name}/{market?.counter.Denom.Name}
        </div>
        <img
          className={`market-arrow ${isOpen ? "rotate" : ""}`}
          src="/trade/images/arrow.svg"
          alt="arr"
        />
      </div>

      <Modal
        isOpen={openCreatePairModal}
        onClose={() => setOpenCreatePairModal(false)}
        title="Create Pair"
        width={640}
      >
        <div className="create-pair-container">
          <div className="create-pair-row">
            <Dropdown
              searchable={true}
              variant={DropdownVariant.OUTLINED}
              items={currencies || []}
              label="Base Token"
              value={baseToken ? baseToken.Denom.Name : market!.base.Denom.Name}
              renderItem={(item: Token) => (
                <div
                  className="create-pair-token"
                  onClick={() => {
                    setBaseToken(item);
                  }}
                >
                  {item.Denom.Name}
                </div>
              )}
            />

            <div className="swap">
              <img
                src="/trade/images/swap.svg"
                alt="swap"
                onClick={() => {
                  setBaseToken(quoteToken);
                  setQuoteToken(baseToken);
                }}
              />
            </div>

            <Dropdown
              searchable={true}
              variant={DropdownVariant.OUTLINED}
              items={currencies || []}
              label="Quote Token"
              value={
                quoteToken ? quoteToken.Denom.Name : market!.counter.Denom.Name
              }
              renderItem={(item: Token) => (
                <div
                  className="create-pair-token"
                  onClick={() => {
                    setQuoteToken(item);
                  }}
                >
                  {item.Denom.Name}
                </div>
              )}
            />
          </div>

          <div className="button-row">
            <Button
              label="Continue"
              width={160}
              disabled={!baseToken || !quoteToken}
              onClick={() => {
                setMarket({
                  base: baseToken!,
                  counter: quoteToken!,
                  pair_symbol: `${baseToken!.Denom.Denom}_${
                    quoteToken!.Denom.Denom
                  }`,
                  reversed_pair_symbol: `${quoteToken!.Denom.Denom}_${
                    baseToken!.Denom.Denom
                  }`,
                });
                setOpenCreatePairModal(false);
              }}
            />
          </div>
        </div>
      </Modal>
    </div>
  );
};

export default MarketSelector;
