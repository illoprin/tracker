import OrangeButton from "@/components/UI/OrangeButton";
import { routes } from "@/router";
import { useNavigate } from "react-router-dom";

function MainPage() {
  const navigate = useNavigate();

  return (
    <>
      <div className="landing-backdrop"></div>

      <div className="flex flex-col w-screen h-screen items-center justify-center">
        <h1 className="text-center text-[16.5rem] text-accent-1 font-thin italic leading-none tracking-[-2.1rem]">
          ТРЕКЕР
        </h1>
        <h1 className="text-center text-lg-1 uppercase font-extralight">
          открыт для нового
        </h1>

        <div className="mt-10">
          <OrangeButton onClick={() => navigate(routes.public.auth.route)}>
            Зарегистрироваться
          </OrangeButton>
        </div>
      </div>
    </>
  );
}

export default MainPage;
