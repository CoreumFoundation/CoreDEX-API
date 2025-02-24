import React from "react";
import DatePicker from "react-datepicker";
import { DatePickerProps } from "react-datepicker";
import "./datetime-picker.scss";

export enum TimeframeType {
  DAY = "day",
  MONTH = "month",
  YEAR = "year",
  QUARTER = "quarter",
  TIME = "time",
}

export type DatetimePickerProps = {
  /** Defines the width of the container and input of the datetime picker*/
  width?: string;
  /** Defines container custom style */
  containerStyle?: React.CSSProperties;
  /** Defines the type of the datetime picker */
  type?: TimeframeType;
  /** Defines the placeholder of the datetime picker */
  placeholder?: string;
  /** Defines the selected date of the datetime picker */
  selectedDate: string | null;
  /** Defines if date range is in the past of the datetime picker*/
  isPast?: boolean;
  /** Defines if date range is in the future of the datetime picker*/
  isFuture?: boolean;
};
type IDatePicker = DatetimePickerProps & DatePickerProps;

export function DatetimePicker({
  containerStyle,
  width,
  type,
  placeholder,
  isPast,
  isFuture,
  selectedDate = null,
  ...rest
}: IDatePicker) {
  return (
    <div
      className={`datetime-container`}
      style={{ ...containerStyle, width: width }}
    >
      <DatePicker
        {...(rest as DatePickerProps)}
        selected={selectedDate ? new Date(selectedDate) : null}
        {...(isPast && { maxDate: new Date() })}
        {...(isFuture && { minDate: new Date() })}
        placeholderText={placeholder}
        showPopperArrow={false}
        yearDropdownItemNumber={15}
        renderCustomHeader={({
          monthDate,
          decreaseMonth,
          increaseMonth,
          prevMonthButtonDisabled,
          nextMonthButtonDisabled,
          decreaseYear,
          increaseYear,
          prevYearButtonDisabled,
          nextYearButtonDisabled,
        }) => (
          <div className="react-datepicker__navigation react-datepicker__navigation--header">
            <div
              aria-label="Previous Month"
              className={
                "react-datepicker__navigation react-datepicker__navigation--previous"
              }
              onClick={() => {
                if (monthDate) !prevMonthButtonDisabled && decreaseMonth();
                else !prevYearButtonDisabled && decreaseYear();
              }}
            >
              <img
                className={
                  "react-datepicker__navigation-icon react-datepicker__navigation-icon--previous"
                }
                style={{ transform: "rotate(90deg)" }}
                src={"/trade/images/arrow.svg"}
              />
            </div>
            <span className="react-datepicker__current-month">
              {monthDate
                ? monthDate.toLocaleString("en-US", {
                    month: "long",
                    year: "numeric",
                  })
                : "Year"}
            </span>
            <div
              aria-label="Next Month"
              className={
                "react-datepicker__navigation react-datepicker__navigation--next"
              }
              onClick={() => {
                if (monthDate) !nextMonthButtonDisabled && increaseMonth();
                else !nextYearButtonDisabled && increaseYear();
              }}
            >
              <img
                className={
                  "react-datepicker__navigation-icon react-datepicker__navigation-icon--next"
                }
                style={{ transform: "rotate(-90deg)" }}
                src={"/trade/images/arrow.svg"}
              />
            </div>
          </div>
        )}
        {...(type === TimeframeType.MONTH && {
          dateFormat: "MM/yyyy",
          showMonthYearPicker: true,
        })}
        {...(type === TimeframeType.YEAR && {
          showYearPicker: true,
          dateFormat: "yyyy",
        })}
        {...(type === TimeframeType.QUARTER && {
          showQuarterYearPicker: true,
          dateFormat: "QQQ, yyyy",
        })}
        {...(type === TimeframeType.TIME && {
          showTimeSelect: true,
          dateFormat: "MM/dd/yyyy h:mm aa",
          filterTime: (time: any) => {
            if (isFuture) {
              const oneHourFromNow = new Date().getTime();
              return oneHourFromNow < new Date(time).getTime();
            } else if (isPast) {
              const oneHourFromNow = new Date().getTime();

              return oneHourFromNow > new Date(time).getTime();
            } else {
              return true;
            }
          },
        })}
        withPortal
      />
    </div>
  );
}
