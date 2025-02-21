import React, { useState } from "react";
import { CopyToClipboard } from "react-copy-to-clipboard";
import { useStore } from "@/state/store";
import "./wallet.scss";
import { resolveCoreumExplorer } from "@/utils";

const Wallet = ({}) => {
  const [isOpen, setIsOpen] = useState(false);
  const { wallet, setLoginModal, pushNotification, network } = useStore();

  const togglewallet = () => setIsOpen((prev) => !prev);

  const walletItems = [
    {
      label: "Copy Address",
      action: () => {
        console.log("copy");
      },
      image: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="21"
          height="21"
          viewBox="0 0 21 21"
          fill="none"
        >
          <mask
            id="path-1-outside-1_186_25904"
            maskUnits="userSpaceOnUse"
            x="2.26953"
            y="1"
            width="16"
            height="19"
            fill="black"
          >
            <rect fill="white" x="2.26953" y="1" width="16" height="19" />
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M8.00227 2C7.26105 2 6.66017 2.60088 6.66017 3.3421V5.39064H4.61163C3.87041 5.39064 3.26953 5.99151 3.26953 6.73274V17.6579C3.26953 18.3991 3.87041 19 4.61163 19H12.9939C13.7351 19 14.336 18.3991 14.336 17.6579V15.6094H16.3845C17.1257 15.6094 17.7266 15.0085 17.7266 14.2673V3.3421C17.7266 2.60088 17.1257 2 16.3845 2H8.00227ZM14.336 14.7146H16.3845C16.6316 14.7146 16.8319 14.5143 16.8319 14.2673V3.3421C16.8319 3.09503 16.6316 2.89473 16.3845 2.89473H8.00227C7.75519 2.89473 7.5549 3.09503 7.5549 3.3421V5.39064H12.9939C13.7351 5.39064 14.336 5.99152 14.336 6.73274V14.7146ZM4.16426 6.73274C4.16426 6.48566 4.36456 6.28537 4.61163 6.28537H12.9939C13.2409 6.28537 13.4412 6.48566 13.4412 6.73274V17.6579C13.4412 17.905 13.2409 18.1053 12.9939 18.1053H4.61163C4.36456 18.1053 4.16426 17.905 4.16426 17.6579V6.73274Z"
            />
          </mask>
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M8.00227 2C7.26105 2 6.66017 2.60088 6.66017 3.3421V5.39064H4.61163C3.87041 5.39064 3.26953 5.99151 3.26953 6.73274V17.6579C3.26953 18.3991 3.87041 19 4.61163 19H12.9939C13.7351 19 14.336 18.3991 14.336 17.6579V15.6094H16.3845C17.1257 15.6094 17.7266 15.0085 17.7266 14.2673V3.3421C17.7266 2.60088 17.1257 2 16.3845 2H8.00227ZM14.336 14.7146H16.3845C16.6316 14.7146 16.8319 14.5143 16.8319 14.2673V3.3421C16.8319 3.09503 16.6316 2.89473 16.3845 2.89473H8.00227C7.75519 2.89473 7.5549 3.09503 7.5549 3.3421V5.39064H12.9939C13.7351 5.39064 14.336 5.99152 14.336 6.73274V14.7146ZM4.16426 6.73274C4.16426 6.48566 4.36456 6.28537 4.61163 6.28537H12.9939C13.2409 6.28537 13.4412 6.48566 13.4412 6.73274V17.6579C13.4412 17.905 13.2409 18.1053 12.9939 18.1053H4.61163C4.36456 18.1053 4.16426 17.905 4.16426 17.6579V6.73274Z"
            fill="#5E6773"
          />
          <path
            d="M6.66017 5.39064V5.59064H6.86017V5.39064H6.66017ZM14.336 15.6094V15.4094H14.136V15.6094H14.336ZM14.336 14.7146H14.136V14.9146H14.336V14.7146ZM7.5549 5.39064H7.3549V5.59064H7.5549V5.39064ZM6.86017 3.3421C6.86017 2.71134 7.3715 2.2 8.00227 2.2V1.8C7.15059 1.8 6.46017 2.49042 6.46017 3.3421H6.86017ZM6.86017 5.39064V3.3421H6.46017V5.39064H6.86017ZM4.61163 5.59064H6.66017V5.19064H4.61163V5.59064ZM3.46953 6.73274C3.46953 6.10197 3.98087 5.59064 4.61163 5.59064V5.19064C3.75995 5.19064 3.06953 5.88106 3.06953 6.73274H3.46953ZM3.46953 17.6579V6.73274H3.06953V17.6579H3.46953ZM4.61163 18.8C3.98087 18.8 3.46953 18.2887 3.46953 17.6579H3.06953C3.06953 18.5096 3.75995 19.2 4.61163 19.2V18.8ZM12.9939 18.8H4.61163V19.2H12.9939V18.8ZM14.136 17.6579C14.136 18.2887 13.6246 18.8 12.9939 18.8V19.2C13.8455 19.2 14.536 18.5096 14.536 17.6579H14.136ZM14.136 15.6094V17.6579H14.536V15.6094H14.136ZM16.3845 15.4094H14.336V15.8094H16.3845V15.4094ZM17.5266 14.2673C17.5266 14.898 17.0153 15.4094 16.3845 15.4094V15.8094C17.2362 15.8094 17.9266 15.1189 17.9266 14.2673H17.5266ZM17.5266 3.3421V14.2673H17.9266V3.3421H17.5266ZM16.3845 2.2C17.0153 2.2 17.5266 2.71134 17.5266 3.3421H17.9266C17.9266 2.49042 17.2362 1.8 16.3845 1.8V2.2ZM8.00227 2.2H16.3845V1.8H8.00227V2.2ZM14.336 14.9146H16.3845V14.5146H14.336V14.9146ZM16.3845 14.9146C16.742 14.9146 17.0319 14.6248 17.0319 14.2673H16.6319C16.6319 14.4039 16.5211 14.5146 16.3845 14.5146V14.9146ZM17.0319 14.2673V3.3421H16.6319V14.2673H17.0319ZM17.0319 3.3421C17.0319 2.98457 16.742 2.69473 16.3845 2.69473V3.09473C16.5211 3.09473 16.6319 3.20548 16.6319 3.3421H17.0319ZM16.3845 2.69473H8.00227V3.09473H16.3845V2.69473ZM8.00227 2.69473C7.64474 2.69473 7.3549 2.98457 7.3549 3.3421H7.7549C7.7549 3.20548 7.86565 3.09473 8.00227 3.09473V2.69473ZM7.3549 3.3421V5.39064H7.7549V3.3421H7.3549ZM7.5549 5.59064H12.9939V5.19064H7.5549V5.59064ZM12.9939 5.59064C13.6246 5.59064 14.136 6.10197 14.136 6.73274H14.536C14.536 5.88106 13.8455 5.19064 12.9939 5.19064V5.59064ZM14.136 6.73274V14.7146H14.536V6.73274H14.136ZM4.61163 6.08537C4.2541 6.08537 3.96426 6.37521 3.96426 6.73274H4.36426C4.36426 6.59612 4.47501 6.48537 4.61163 6.48537V6.08537ZM12.9939 6.08537H4.61163V6.48537H12.9939V6.08537ZM13.6412 6.73274C13.6412 6.37521 13.3514 6.08537 12.9939 6.08537V6.48537C13.1305 6.48537 13.2412 6.59612 13.2412 6.73274H13.6412ZM13.6412 17.6579V6.73274H13.2412V17.6579H13.6412ZM12.9939 18.3053C13.3514 18.3053 13.6412 18.0154 13.6412 17.6579H13.2412C13.2412 17.7945 13.1305 17.9053 12.9939 17.9053V18.3053ZM4.61163 18.3053H12.9939V17.9053H4.61163V18.3053ZM3.96426 17.6579C3.96426 18.0154 4.2541 18.3053 4.61163 18.3053V17.9053C4.47501 17.9053 4.36426 17.7945 4.36426 17.6579H3.96426ZM3.96426 6.73274V17.6579H4.36426V6.73274H3.96426Z"
            fill="#5E6773"
            mask="url(#path-1-outside-1_186_25904)"
          />
        </svg>
      ),
    },
    {
      label: "Switch Wallet",
      action: () => {
        setLoginModal(true);
      },
      image: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="21"
          height="21"
          viewBox="0 0 21 21"
          fill="none"
        >
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M6.91616 5.58207L9.74083 5.58208L7.91951 3.76075C7.67232 3.51356 7.67232 3.14757 7.91951 2.90039C8.1667 2.6532 8.53268 2.6532 8.77987 2.90039L11.6469 5.76744C11.8941 6.01463 11.8941 6.3806 11.6469 6.62779C11.5812 6.69349 11.5195 6.73875 11.4534 6.76817C11.3877 6.79738 11.3125 6.81322 11.2168 6.81322H6.91616C4.85334 6.81322 3.23115 8.43541 3.23115 10.4982C3.23115 12.5611 4.85334 14.1833 6.91616 14.1833C7.1082 14.1833 7.26036 14.2468 7.36426 14.3507C7.46817 14.4546 7.53174 14.6068 7.53174 14.7988C7.53174 14.9909 7.46817 15.143 7.36426 15.2469C7.26036 15.3508 7.1082 15.4144 6.91616 15.4144C4.17667 15.4144 2 13.2377 2 10.4982C2 7.75874 4.17667 5.58207 6.91616 5.58207ZM14.0838 15.4142L11.2591 15.4142L13.0805 17.2356C13.3277 17.4827 13.3277 17.8487 13.0805 18.0959C12.8333 18.3431 12.4673 18.3431 12.2201 18.0959L9.35304 15.2289C9.10586 14.9817 9.10586 14.6157 9.35304 14.3685C9.41875 14.3028 9.48043 14.2576 9.54657 14.2281C9.61228 14.1989 9.6875 14.1831 9.78318 14.1831H14.0838C16.1466 14.1831 17.7688 12.5609 17.7688 10.4981C17.7688 8.43525 16.1466 6.81306 14.0838 6.81306C13.8918 6.81306 13.7396 6.74949 13.6357 6.64559C13.5318 6.54168 13.4682 6.38952 13.4682 6.19748C13.4682 6.00545 13.5318 5.85328 13.6357 5.74938C13.7396 5.64548 13.8918 5.58191 14.0838 5.58191C16.8233 5.58191 19 7.75858 19 10.4981C19 13.2376 16.8233 15.4142 14.0838 15.4142Z"
            fill="#5E6773"
          />
        </svg>
      ),
    },
    {
      label: "Open Explorer",
      action: () => {
        window.open(
          `${resolveCoreumExplorer(network)}/accounts/${wallet.address}`
        );
      },
      image: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="21"
          height="21"
          viewBox="0 0 21 21"
          fill="none"
        >
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M2 5.941C2.00498 5.83929 2.03774 5.74086 2.09471 5.65645C2.15166 5.57191 2.23083 5.50453 2.32345 5.46193L10.2992 2.04202C10.43 1.98599 10.578 1.98599 10.7087 2.04202L18.6849 5.46208C18.7777 5.50308 18.8566 5.56994 18.9126 5.65471C18.9685 5.73949 18.9988 5.83885 19 5.94037L19 15.0602C18.9999 15.162 18.9698 15.2617 18.9135 15.3465C18.8572 15.4314 18.7771 15.4981 18.6834 15.5378L10.703 18.958C10.5728 19.014 10.4253 19.014 10.2952 18.958L2.31494 15.5378C2.22147 15.4978 2.14189 15.4312 2.08584 15.3463C2.02981 15.2614 2 15.162 2 15.0602V5.941ZM3.03967 14.7178L9.9805 17.6924V9.70242L3.03967 6.72761V14.7178ZM3.83948 5.93982L10.5002 8.79507L17.161 5.93982L10.5002 3.0851L3.83948 5.93982ZM11.0199 17.6924L17.9608 14.7178V6.72761L11.0199 9.70229V17.6924Z"
            fill="#5E6773"
          />
        </svg>
      ),
    },
    {
      label: "Disconnect",
      action: async () => {
        try {
          localStorage.removeItem("wallet");
          window.location.reload();
        } catch (e) {
          console.log("E_DISCONNECT =>", e);
        }
      },
      image: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="21"
          height="21"
          viewBox="0 0 21 21"
          fill="none"
        >
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M4.61205 3.27658C4.24327 2.90781 3.64536 2.90781 3.27658 3.27658C2.90781 3.64536 2.90781 4.24327 3.27658 4.61205L9.16453 10.5L3.27658 16.3879C2.90781 16.7567 2.90781 17.3546 3.27658 17.7234C3.64536 18.0922 4.24327 18.0922 4.61205 17.7234L10.5 11.8355L16.3877 17.7232C16.7565 18.092 17.3544 18.092 17.7232 17.7232C18.092 17.3544 18.092 16.7565 17.7232 16.3877L11.8355 10.5L17.7232 4.61228C18.092 4.2435 18.092 3.64559 17.7232 3.27681C17.3544 2.90803 16.7565 2.90803 16.3877 3.27681L10.5 9.16453L4.61205 3.27658Z"
            fill="#5E6773"
          />
        </svg>
      ),
    },
  ];

  return (
    <div className={`wallet wallet`}>
      <button className="wallet-label" onClick={togglewallet}>
        <div className="wallet-label-selected">{wallet?.address}</div>
        <img
          className={`wallet-arrow ${isOpen ? "rotate" : ""}`}
          src="/trade/images/arrow.svg"
          alt="arr"
        />
      </button>

      <div className={`wallet-list ${isOpen ? "open" : ""}`}>
        <ul className="wallet-list-content">
          {walletItems.map((item, index) => (
            <li key={index}>
              {item.label === "Copy Address" ? (
                <CopyToClipboard
                  text={wallet?.address || ""}
                  onCopy={() => {
                    console.log("copied");
                    pushNotification({
                      message: "Address copied to clipboard",
                      type: "success",
                    });
                  }}
                >
                  <div
                    className="wallet-item"
                    style={
                      {
                        "--default-color": "#5e6773",
                        "--hover-color": "#eee",
                      } as React.CSSProperties
                    }
                    onClick={() => {
                      try {
                        item.action();
                      } catch {
                        console.log("error");
                      }
                    }}
                  >
                    <div className="image">{item.image}</div>
                    <div className="label">{item.label}</div>
                  </div>
                </CopyToClipboard>
              ) : (
                <div
                  className="wallet-item"
                  style={
                    {
                      "--default-color": "#5e6773",
                      "--hover-color": "#eee",
                    } as React.CSSProperties
                  }
                  onClick={() => {
                    try {
                      item.action();
                    } catch {
                      console.log("error");
                    }
                  }}
                >
                  <div className="image">{item.image}</div>
                  <div className="label">{item.label}</div>
                </div>
              )}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default Wallet;
