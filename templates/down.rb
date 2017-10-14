# ExMachina

# Plays chords when a pod is going down ;-(

# Author Fernand Galiana
debug = false
use_bpm 120
use_random_seed 100
load_sample :sn_dub

n = [:f6, :g6, :e6].choose
with_fx :reverb, mix: 0.6, room: 0.8 do
  with_fx :echo, room: 0.8, decay: 8, phase: 1, mix: 0.4 do
    with_synth :blade do
      with_transpose -12 do
        in_thread do
          {{ .Count }}.times do
            play n, amp: 6, attack: 0.6, release: 0.8, detune: rrand(0, 0.1), cutoff: rrand(80, 120)
            sleep 2
          end
        end
      end
    end

    sleep 2
    sample :sn_dub, beat_stretch: 2.3, decay: 1.5
    with_synth :tech_saws do
      play chord(n, :m7), amp: 5, release: 0.8
    end
    sleep 2
  end
end