import HeartIcon from "@/components/shared/icons/HeartIcon";

function TrackItem({ track }) {
  return (
    <div className="flex items-center space-x-4 p-3 rounded-md hover:bg-background-secondary transition-colors group">
      <div className="w-12 h-12 bg-background-light rounded-md flex items-center justify-center">
        <div className="w-8 h-8 bg-accent-1 rounded-sm"></div>
      </div>
      <div className="flex-1">
        <h3 className="text-text font-medium">{track.title}</h3>
        <p className="text-secondary text-sm">{track.artist}</p>
      </div>
      <span className="text-secondary text-sm">{track.duration}</span>
      <button className="opacity-0 group-hover:opacity-100 transition-opacity">
        <HeartIcon className="w-5 h-5 text-secondary hover:text-accent-1" />
      </button>
    </div>
  );
}

export default TrackItem;
