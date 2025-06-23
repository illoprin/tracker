import LockIcon from '@/components/shared/icons/LockIcon';
import UserIcon from '@/components/shared/icons/UserIcon';
import BaseHeader from '@/components/UI/BaseHeader';
import DarkInput from '@/components/UI/DarkInput';
import OrangeButton from '@/components/UI/OrangeButton';
import { routes } from '@/router';
import { Link } from 'react-router-dom';

function LoginForm() {
  return (
    <>
      <BaseHeader>С возвращением!</BaseHeader>


      <div className="flex gap-3 flex-col">

        {/* Login */}
        <DarkInput placeholder="Логин" type="text">
          <UserIcon className="absolute top-[25%] left-3 stroke-secondary" />
        </DarkInput>

        {/* Password */}
        <DarkInput placeholder="Пароль" type="password">
          <LockIcon className="absolute top-[25%] left-3 stroke-secondary" />
        </DarkInput>

        <OrangeButton>Войти</OrangeButton>
        
      </div>

      <div className="mt-14 w-[100%] text-center text-md">
        <span className="text-secondary">Нет аккаунта?</span>
        <Link
          to={{ hash: routes.public.auth.hashes.reg }}
          className="underline text-accent-1 text-xl ml-2"
        >
          Зарегистрироваться
        </Link>
      </div>
    </>
  );
}

export default LoginForm;
