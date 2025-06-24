import DownChevron from "@/components/shared/icons/DownChevron";
import PlusIcon from "@/components/shared/icons/PlusIcon";
import SearchIcon from "@/components/shared/icons/SearchIcon";
import DarkInput from "@/components/UI/DarkInput";
import GrayButton from "@/components/UI/GrayButton";

function NavBar() {
  

  return (
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
        <SearchIcon className="absolute top-[25%] left-3 stroke-secondary" />
      </DarkInput>

      <div className="flex items-center space-x-4">
        <GrayButton>
          <PlusIcon className="block float-start stroke-secondary" />
          <DownChevron className="block float-start stroke-secondary" />
        </GrayButton>
        <button className="flex items-center space-x-2 text-secondary hover:text-text transition-colors">
          <span>Логин</span>
        </button>
      </div>
    </header>
  );
}

export default NavBar;
