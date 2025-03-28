import { useEffect, useRef, useState } from "react";
import Button from "../Button";
import Modal from "../Modal";
import Dropdown, { DropdownVariant } from "../Dropdown";
import { useStore } from "@/state/store";
import { Token } from "@/types/market";
import { getCurrencies } from "@/services/api";
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
          handleDevcore(currs);
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

  // rename devcore to core for display purposes + add to top of list
  const handleDevcore = (currs: Token[]) => {
    const devcore = currs.find((curr) => curr.Denom.Name === "devcore");
    if (devcore) {
      devcore.Denom.Name = "Core";
      devcore.Denom.Currency = "Core";
      currs = currs.filter((curr) => curr.Denom.Name !== "Core");
      currs.unshift(devcore!);
    }

    setCurrencies(currs);
  };

  return (
    <div className="market-selector-container" ref={ref}>
      <div
        className="market-label"
        onClick={() => {
          setOpenCreatePairModal(true);
        }}
      >
        <div className="market-label-selected">
          {market?.base.Denom.Name?.toUpperCase()}/
          {market?.counter.Denom.Name?.toUpperCase()}
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
              onClick={(item) => {
                setBaseToken(item);
              }}
              renderItem={(item: Token) => (
                <div className="create-pair-token">
                  <p className="create-pair-name">{item.Denom.Name}</p>
                  <p className="create-pair-issuer">{item.Denom.Issuer}</p>
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
                quoteToken ? quoteToken.Denom.Name : market!.counter.Denom?.Name
              }
              onClick={(item) => {
                setQuoteToken(item);
              }}
              renderItem={(item: Token) => (
                <div className="create-pair-token">
                  <p className="create-pair-name">{item.Denom.Name}</p>
                  <p className="create-pair-issuer">{item.Denom.Issuer}</p>
                </div>
              )}
            />
          </div>

          <div className="button-row">
            <Button
              label="Confirm"
              width={160}
              disabled={
                !baseToken ||
                !quoteToken ||
                (baseToken.Denom.Currency === quoteToken.Denom.Currency &&
                  baseToken.Denom.Issuer === quoteToken.Denom.Issuer)
              }
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
