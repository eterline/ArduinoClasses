#include "PCF8574.h"

class ButtonPCF8574 {
  private:
    PCF8574*  _pcf_ptr;
    uint8_t   pcf_btn_pin;

    bool      pullUp;
    bool currentState = false;
    bool lastState = false;
    bool clickFlag = false;

    unsigned long lastDebounceTime = 0;
    const unsigned long debounceDelay = 30;

    unsigned long clickTime = 0;
    const unsigned long clickResetDelay = 500;

  public:
    ButtonPCF8574(PCF8574* pcf, uint8_t btn_pin, bool pull_upped = true) {
      _pcf_ptr = pcf;
      pcf_btn_pin = btn_pin;
      pullUp = pull_upped;
      _pcf_ptr->pinMode(pcf_btn_pin, INPUT);
    }

    void doTick() {
      bool rawState = _pcf_ptr->digitalRead(pcf_btn_pin);
      bool reading = pullUp ? !rawState : rawState;
      unsigned long now = millis();

      if (reading != lastState) {
        lastDebounceTime = now;
      }

      if (currentState && (now - clickTime > clickResetDelay)) {
        clickTime = now;
        currentState = false;
      }

      if ((now - lastDebounceTime) > debounceDelay) {
        if (reading != currentState) {
          currentState = reading;

          if (currentState) {
            clickFlag = true;
            clickTime = now;
          }
        }
      }

      lastState = reading;
    }

    bool wasClick() {
      if (clickFlag) {
        clickFlag = false;
        return true;
      }
      return false;
    }
};