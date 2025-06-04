import React, { useState, useEffect, useRef } from "react";
import "./dropdown.scss";

interface DropdownProps<T> {
  value: string | undefined;
  items: T[];
  renderItem: (item: T, index: number) => React.ReactNode;
  getvalue?: (item: T) => string;
  variant?: DropdownVariantType;
  image?: string;
  label?: string;
  searchable?: boolean;
  onClick?: (item: T) => void;
}

export const DropdownVariant = {
  DEFAULT: "default",
  OUTLINED: "outlined",
  NETWORK: "network",
};

export type DropdownVariantType =
  (typeof DropdownVariant)[keyof typeof DropdownVariant];

const Dropdown = <T,>({
  value,
  items,
  getvalue,
  renderItem,
  variant = "default",
  image,
  label,
  searchable = false,
  onClick,
}: DropdownProps<T>) => {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedvalue, setSelectedvalue] = useState<string | undefined>(value);
  const [searchQuery, setSearchQuery] = useState<string>(""); // State for search query
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setSelectedvalue(value);
  }, [value]);

  useEffect(() => {
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const filteredItems = searchable
    ? items.filter((item: any) => {
        const itemValue = getvalue ? getvalue(item) : item.Denom?.Name;
        return itemValue.toLowerCase().includes(searchQuery.toLowerCase());
      })
    : items;

  const toggleDropdown = () => setIsOpen((prev) => !prev);
  const closeDropdown = () => {
    setSearchQuery("");
    setIsOpen(false);
  };

  const handleItemClick = (item: any) => {
    if (item.hasOwnProperty("Denom")) {
      setSelectedvalue(item.Denom.Name);
    } else {
      setSelectedvalue(getvalue ? getvalue(item) : (item as unknown as string));
    }
    if (onClick) {
      onClick(item);
    }
    closeDropdown();
  };

  const handleClickOutside = (event: MouseEvent) => {
    if (
      dropdownRef.current &&
      !dropdownRef.current.contains(event.target as Node)
    ) {
      closeDropdown();
    }
  };

  return (
    <div className="dropdown-container" ref={dropdownRef}>
      {label && <div className="dropdown-label">{label}</div>}
      <div className={`dropdown dropdown-${variant}`}>
        <button
          className={`dropdown-value ${isOpen ? "active" : ""}`}
          onClick={toggleDropdown}
        >
          <div className="dropdown-value-left">
            {image && (
              <img
                src={image}
                alt={value}
                className="dropdown-selected-image"
              />
            )}
            <div className="dropdown-value-selected">
              {selectedvalue
                ? selectedvalue.length > 30
                  ? selectedvalue.slice(0, 30) + "..."
                  : selectedvalue
                : ""}
            </div>
          </div>

          <img
            className={`dropdown-arrow ${isOpen ? "rotate" : ""}`}
            src="/trade/images/arrow.svg"
            alt="arr"
          />
        </button>

        <div className={`dropdown-list ${variant} ${isOpen ? "open" : ""}`}>
          {searchable && (
            <div className="dropdown-search">
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search..."
                className="dropdown-search-input"
              />
            </div>
          )}

          <ul className="dropdown-list-content">
            {filteredItems.map((item, index) => (
              <li
                key={index}
                className="dropdown-item"
                onClick={() => handleItemClick(item)}
                style={{ display: "flex", justifyContent: "space-between" }}
              >
                {renderItem(item, index)}
                {getvalue
                  ? getvalue(item) === selectedvalue
                  : (item as unknown as string) === selectedvalue && (
                      <img
                        style={{
                          width: "20px",
                          height: "20px",
                          marginRight: "20px",
                        }}
                        src="/trade/images/check.svg"
                        alt="selected"
                      />
                    )}
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  );
};

export default Dropdown;
