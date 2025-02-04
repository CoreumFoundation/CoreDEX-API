import Orderbook from "@/components/Orderbook";
import TradingView from "@/components/TradingView";
import OrderActions from "@/components/OrderActions";
import "./home.scss";
import Header from "@/components/Header";
import { useState } from "react";
import { OrderbookAction } from "@/types/market";
import Modal from "@/components/Modal";
import { useStore } from "@/state/store";
import LoginSelection from "@/components/LoginSelection";
import ExchangeHistory from "@/components/ExchangeHistory";
import OrderHistory from "@/components/OrderHistory";
import { useWindowSize } from "react-use";

const Home = () => {
  const [orderbookAction, setOrderbookAction] = useState<OrderbookAction>();
  const { loginModal, setLoginModal } = useStore();
  const { width } = useWindowSize();

  return (
    <div className="home-container">
      <Header />

      {width > 768 ? (
        <>
          <div className="row">
            <div className="stack">
              <Orderbook setOrderbookAction={setOrderbookAction} />
              <ExchangeHistory />
            </div>
            <TradingView height="100%" />
            <OrderActions orderbookAction={orderbookAction} />
          </div>

          <div className="row">
            <OrderHistory />
          </div>
        </>
      ) : (
        <>
          <div className="row">
            <TradingView height="100%" />
            <OrderActions orderbookAction={orderbookAction} />
          </div>

          <div className="row">
            <Orderbook setOrderbookAction={setOrderbookAction} />
            <ExchangeHistory />
          </div>
          <div className="row">
            <OrderHistory />
          </div>
        </>
      )}

      <Modal
        isOpen={loginModal}
        title="Connect Coreum Wallet"
        onClose={() => setLoginModal(false)}
      >
        <LoginSelection closeModal={() => setLoginModal(false)} />

        <p className="wallet-disclaimer">
          * For devnet, Keplr is only working extension
        </p>
      </Modal>
    </div>
  );
};

export default Home;
