#include "Arduino.h"
#include "hardware/adc.h"
#include "hardware/pwm.h"

// ポテンションメーター２番ピン読み取り用
const int potPin = 26;

// ポテンションメーター初期値
int potValue = 0;

// PWMスライスを取得（２番のやつ）
uint slice_num = pwm_gpio_to_slice_num(2);

void setup()
{

  // シリアル通信を開始 (ボーレート115200)
  // Serial.begin(115200);

  gpio_set_function(2, GPIO_FUNC_PWM); // PWM用に設定
  gpio_set_function(3, GPIO_FUNC_PWM); // PWM用に設定

  // PWM周期を設定
  pwm_set_clkdiv(slice_num, 1.220703125);
  // アナログピンの読み取り値に合わせて12bitMaxでWrapする
  pwm_set_wrap(slice_num, 4096);

  // A,Bともに0Start
  pwm_set_chan_level(slice_num, PWM_CHAN_A, 0); // GPIO2
  pwm_set_chan_level(slice_num, PWM_CHAN_B, 0); // GPIO2+1=GPIO3

  // PWM出力イネーブル
  pwm_set_enabled(slice_num, true);

  analogReadResolution(12);
  pinMode(potPin, INPUT);
}

void loop()
{
  // ポテンショメーターのアナログ値を読み取る (0-4096)
  int potValue = analogRead(potPin);

  // 読み取ったアナログ値をそのまま渡す（wrapにて設定済み）
  pwm_set_chan_level(slice_num, PWM_CHAN_A, potValue);
  pwm_set_chan_level(slice_num, PWM_CHAN_B, potValue);
  // シリアルモニタにポテンショメーターの値とPWM値を出力
  // Serial.print("Pot Value: ");
  // Serial.println(potValue);

  // 少し待機
  delay(10);
}