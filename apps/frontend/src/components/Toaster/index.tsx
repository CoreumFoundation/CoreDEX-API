import { CSSProperties } from "react";
import { ToastContainer, ToastPosition } from "react-toastify";
import "./toaster.scss";

interface ToasterProps {
  customStyle?: CSSProperties;
  position?: ToastPosition;
}

export const Toaster = ({
  customStyle,
  position = "bottom-right",
}: ToasterProps) => {
  return (
    <ToastContainer
      style={{ zIndex: 8000, ...customStyle }}
      autoClose={4000}
      closeOnClick
      hideProgressBar
      pauseOnHover
      theme={"colored"}
      position={position}
      className={`toaster-container `}
      closeButton={({ closeToast }) => (
        <img
          src=""
          style={{
            width: 17,
            marginLeft: 12,
            marginRight: -4,
          }}
          onClick={closeToast}
        />
      )}
    />
  );
};
