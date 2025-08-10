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



  // PWM信号の周波数について
  //
  // clkdiv: クロック分周
  // PWM周波数 : システムクロック / clkdiv * wrap
  //
  // 125,000,000 / 1.220703125 * 4096 = 102,400,000 / 4096 = 25,000Hz(=25kHz)
  //
  // わざわざクロック周波数を25kHzに変更しているのは、一般的にDCモーターをPWM
  // 制御する場合、20kHz以上にすることが推奨されているため。
  // それ以下の場合、コイル鳴きなどの不快な音を発生させる可能性が高まるようだ。
  // ちなみに、それ以下で制御できないわけではない。

  //
  // PWM周期を設定
  //

  // pwm_set_clkdiv (クロック分周)
  // Picoの基本クロック（通常は 125MHz）を「どれくらい遅くするか」を決める設定。
  // 1.220703125 を設定し基本クロックを遅くする
  //
  // 125,000,000 Hz / 1.220703125 = 102,400,000 Hz
  //
  pwm_set_clkdiv(slice_num, 1.220703125);
  // pwm_set_wrap (ラップ値)
  // 遅くしたPicoの基本クロックを「いくつまで数えたら1サイクル（ON/OFFの1回）とするか」を決める
  // 4096 は、102.4MHzのクロックを4096回カウントしたら、PWMの波一つ分ということ
  // アナログピンの読み取り値に合わせて12bitMaxでWrapするため 4096 を指定している
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