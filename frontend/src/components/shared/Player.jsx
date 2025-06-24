import BackwardStepIcon from "@/components/shared/icons/BackwardStepIcon";
import ForwardStepIcon from "@/components/shared/icons/ForwardStepIcon";
import HeartIcon from "@/components/shared/icons/HeartIcon";
import PauseIcon from "@/components/shared/icons/PauseIcon";
import PlayIcon from "@/components/shared/icons/PlayIcon";
import QueueIcon from "@/components/shared/icons/QueueIcon";
import RepeatIcon from "@/components/shared/icons/RepeatIcon";
import ShuffleIcon from "@/components/shared/icons/ShuffleIcon";
import VolumeIcon from "@/components/shared/icons/VolumeIcon";
import WaveIcon from "@/components/shared/icons/WaveIcon";
import { useState } from "react";

function Player({}) {
  const [isPlaying, setIsPlaying] = useState(true);
  const [currentTime, setCurrentTime] = useState(37);
  const [duration, setDuration] = useState(179);
  const [volume, setVolume] = useState(80);

  const formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="fixed bottom-0 left-0 right-0 bg-background-secondary border-t border-background-light p-4">
      <div className="flex items-center justify-between">
        {/* Current Track Info */}
        <div className="flex items-center space-x-4 flex-1">
          <div className="w-12 h-12 bg-pink-600 rounded-md flex items-center justify-center">
            <div className="text-white font-bold text-xs">АВ</div>
          </div>
          <div>
            <h4 className="text-text font-medium text-sm">Сердце</h4>
            <p className="text-secondary text-xs">Алексей Вишня</p>
            <p className="text-secondary text-xs">Сердце</p>
          </div>
        </div>

        {/* Player Controls */}
        <div className="flex-2 flex flex-col items-center space-y-2">
          <div className="flex items-center space-x-4">
            <button className="text-secondary hover:text-text transition-colors">
              <ShuffleIcon />
            </button>
            <button className="text-secondary hover:text-text transition-colors">
              <RepeatIcon />
            </button>
            <button
              className="bg-accent-1 text-background p-2 rounded-full hover:bg-accent-2 transition-colors"
              onClick={() => setIsPlaying(!isPlaying)}
            >
              {isPlaying ? <PauseIcon /> : <PlayIcon />}
            </button>
            <button className="text-secondary hover:text-text transition-colors">
              <BackwardStepIcon />
            </button>
            <button className="text-secondary hover:text-text transition-colors">
              <ForwardStepIcon />
            </button>
          </div>

          {/* Progress Bar */}
          <div className="flex items-center space-x-3 w-80">
            <span className="text-secondary text-xs">
              0:{currentTime.toString().padStart(2, '0')}
            </span>
            <div className="flex-1 h-1 bg-background-light rounded-full">
              <div
                className="h-full bg-accent-1 rounded-full relative"
                style={{
                  width: `${(currentTime / duration) * 100}%`,
                }}
              >
                <div className="absolute right-0 top-1/2 transform translate-x-1/2 -translate-y-1/2 w-3 h-3 bg-accent-1 rounded-full"></div>
              </div>
            </div>
            <span className="text-secondary text-xs">
              {formatTime(duration)}
            </span>
          </div>
        </div>

        {/* Right Controls */}
        <div className="flex items-center space-x-4 flex-1 justify-end">
          <button className="text-secondary hover:text-text transition-colors">
            <HeartIcon />
          </button>
          <button className="text-secondary hover:text-text transition-colors">
            <WaveIcon />
          </button>
          <button className="text-secondary hover:text-text transition-colors">
            <QueueIcon />
          </button>
          <div className="flex items-center space-x-2">
            <VolumeIcon />
            <div className="w-20 h-1 bg-background-light rounded-full">
              <div
                className="h-full bg-accent-1 rounded-full relative"
                style={{ width: `${volume}%` }}
              >
                <div className="absolute right-0 top-1/2 transform translate-x-1/2 -translate-y-1/2 w-3 h-3 bg-accent-1 rounded-full"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Player;
