import BackwardStepIcon from '@/components/shared/icons/BackwardStepIcon';
import DownChevron from '@/components/shared/icons/DownChevron';
import ForwardStepIcon from '@/components/shared/icons/ForwardStepIcon';
import HeartIcon from '@/components/shared/icons/HeartIcon';
import PauseIcon from '@/components/shared/icons/PauseIcon';
import PlayIcon from '@/components/shared/icons/PlayIcon';
import PlusIcon from '@/components/shared/icons/PlusIcon';
import QueueIcon from '@/components/shared/icons/QueueIcon';
import RepeatIcon from '@/components/shared/icons/RepeatIcon';
import SearchIcon from '@/components/shared/icons/SearchIcon';
import ShuffleIcon from '@/components/shared/icons/ShuffleIcon';
import VolumeIcon from '@/components/shared/icons/VolumeIcon';
import WaveIcon from '@/components/shared/icons/WaveIcon';
import DarkInput from '@/components/UI/DarkInput';
import { useState } from 'react';

function HomePage() {
  const [isPlaying, setIsPlaying] = useState(true);
  const [currentTime, setCurrentTime] = useState(37);
  const [duration, setDuration] = useState(179);
  const [volume, setVolume] = useState(80);

  const recommendations = [
    {
      title: 'People Need Love',
      artist: 'ABBA',
      duration: '2:47',
      cover: '/api/placeholder/60/60',
    },
    {
      title: 'Dance With Me (Remastered 2021)',
      artist: 'Alphaville',
      duration: '4:00',
      cover: '/api/placeholder/60/60',
    },
    {
      title: 'Take A Chance On Me',
      artist: 'ABBA',
      duration: '4:05',
      cover: '/api/placeholder/60/60',
    },
    {
      title: 'Мы не ждали перемен',
      artist: 'Александр Градский',
      duration: '2:58',
      cover: '/api/placeholder/60/60',
    },
    {
      title: 'Расчёска',
      artist: 'Алексей Вишня',
      duration: '3:35',
      cover: '/api/placeholder/60/60',
    },
    {
      title: 'Sounds Like a Melody',
      artist: 'Alphaville',
      duration: '2:34',
      cover: '/api/placeholder/60/60',
    },
  ];

  const listenNext = [
    {
      title: 'Сердце',
      artist: 'Алексей Вишня',
      cover: '/api/placeholder/150/150',
      color: 'bg-pink-600',
    },
    {
      title: 'Life in Unknown City',
      artist: 'FM Skyline',
      cover: '/api/placeholder/150/150',
      color: 'bg-blue-500',
    },
    {
      title: 'To Be Or Not',
      artist: 'Saint Preux',
      cover: '/api/placeholder/150/150',
      color: 'bg-gray-400',
    },
  ];

  const formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <>
      {/* Header */}
      <header className="flex items-center justify-between p-6 border-b border-background-light">
        <div className="flex items-center space-x-8">
          <div className="text-accent-1 text-2xl font-bold">TP</div>
          <nav className="flex space-x-6">
            <a
              href="#"
              className="text-text hover:text-accent-1 transition-colors"
            >
              Домой
            </a>
            <a
              href="#"
              className="text-secondary hover:text-text transition-colors"
            >
              Мой выбор
            </a>
          </nav>
        </div>

        <DarkInput type="text" placeholder="Что хочешь включить?">
          <SearchIcon className="absolute top-[25%] left-3 stroke-secondary"/>
        </DarkInput>

        <div className="flex items-center space-x-4">
          <button className="bg-background-secondary p-2 rounded-md hover:bg-background-light transition-colors">
            <PlusIcon />
          </button>
          <button className="flex items-center space-x-2 text-secondary hover:text-text transition-colors">
            <span>Логин</span>
            <DownChevron />
          </button>
        </div>
      </header>

      <div className="flex">
        {/* Main Content */}
        <main className="flex-1 p-6">
          {/* Recommendations Section */}
          <section className="mb-8">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-semibold">
                Рекомендуем послушать
              </h2>
              <button className="text-secondary hover:text-accent-1 transition-colors text-sm">
                Показать все
              </button>
            </div>

            <div className="space-y-3">
              {recommendations.map((track, index) => (
                <div
                  key={index}
                  className="flex items-center space-x-4 p-3 rounded-md hover:bg-background-secondary transition-colors group"
                >
                  <div className="w-12 h-12 bg-background-light rounded-md flex items-center justify-center">
                    <div className="w-8 h-8 bg-accent-1 rounded-sm"></div>
                  </div>
                  <div className="flex-1">
                    <h3 className="text-text font-medium">
                      {track.title}
                    </h3>
                    <p className="text-secondary text-sm">
                      {track.artist}
                    </p>
                  </div>
                  <span className="text-secondary text-sm">
                    {track.duration}
                  </span>
                  <button className="opacity-0 group-hover:opacity-100 transition-opacity">
                    <HeartIcon className="w-5 h-5 text-secondary hover:text-accent-1" />
                  </button>
                </div>
              ))}
            </div>
          </section>
        </main>

        {/* Sidebar */}
        <aside className="w-80 p-6 border-l border-background-light">
          <div className="mb-6">
            <h2 className="text-lg font-semibold mb-4">
              Слушаем дальше?
            </h2>
            <div className="grid grid-cols-1 gap-4">
              {listenNext.map((item, index) => (
                <div
                  key={index}
                  className="relative group cursor-pointer"
                >
                  <div
                    className={`w-full h-32 ${item.color} rounded-md flex items-center justify-center relative overflow-hidden`}
                  >
                    <div className="absolute inset-0 bg-black bg-opacity-20"></div>
                    <div className="relative text-center">
                      <h3 className="text-white font-semibold text-sm">
                        {item.title}
                      </h3>
                      <p className="text-white text-xs opacity-80">
                        {item.artist}
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </aside>
      </div>

      {/* Bottom Player */}
      <div className="fixed bottom-0 left-0 right-0 bg-background-secondary border-t border-background-light p-4">
        <div className="flex items-center justify-between">
          {/* Current Track Info */}
          <div className="flex items-center space-x-4 flex-1">
            <div className="w-12 h-12 bg-pink-600 rounded-md flex items-center justify-center">
              <div className="text-white font-bold text-xs">АВ</div>
            </div>
            <div>
              <h4 className="text-text font-medium text-sm">
                Сердце
              </h4>
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
    </>
  );
}

export default HomePage;
