import Player from '@/components/shared/Player';
import TrackItem from '@/components/shared/TrackItem';
import NavBar from '@/layouts/NavBar';

function HomePage() {
  

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

  

  return (
    <>
      {/* Header */}
      <NavBar />

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
                <TrackItem key={index} track={track} />
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
      <Player />
    </>
  );
}

export default HomePage;
