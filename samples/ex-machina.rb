# ExMachina

# Coz your K8s cluster got to have more CowBell, Baby!

# Coded by Fernand Galiana

use_debug false
use_bpm 90
use_random_seed 100
load_samples :drum_tom_lo_soft, :drum_heavy_kick, :drum_cymbal_closed
load_sample :ambi_lunar_land

8.times do
  sample :drum_tom_lo_hard, amp: (line 0, 3, steps: 8).tick
  sleep 0.3
end

live_loop :drums do
  n = [:e2, :e2, :f3, :f3, :g5].choose
  sample :drum_tom_lo_soft, amp: 2
  sleep 0.35
  sample :drum_heavy_kick, amp: 2, rate: 0.8
  sleep 0.05
  sample :drum_cymbal_closed, amp: 1
  sleep 0.3
  with_synth :tech_saws do
    play chord(n, :m7), amp: 0.9, release: 0.8
  end
end

live_loop :bar, auto_cue: false do
  if rand < 0.25
    sample :ambi_lunar_land
    puts :comet_landing
  end
  sleep 8
end