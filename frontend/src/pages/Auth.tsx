import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { useToast } from '@/hooks/use-toast';
import { Loader2 } from 'lucide-react';

const authSchema = z.object({
  name: z.string().min(3, 'Имя должно содержать минимум 3 символа').max(50, 'Имя слишком длинное'),
  password: z.string().min(6, 'Пароль должен содержать минимум 6 символов').max(100, 'Пароль слишком длинный'),
});

type AuthFormValues = z.infer<typeof authSchema>;

export default function Auth() {
  const [isSignUp, setIsSignUp] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { signIn, signUp } = useAuth();
  const navigate = useNavigate();
  const { toast } = useToast();

  const form = useForm<AuthFormValues>({
    resolver: zodResolver(authSchema),
    defaultValues: {
      name: '',
      password: '',
    },
  });

  const onSubmit = async (values: AuthFormValues) => {
    setIsSubmitting(true);
    
    const result = isSignUp 
      ? await signUp(values.name, values.password)
      : await signIn(values.name, values.password);

    setIsSubmitting(false);

    if (result.error) {
      toast({
        variant: 'destructive',
        title: isSignUp ? 'Ошибка регистрации' : 'Ошибка входа',
        description: result.error,
      });
      return;
    }

    toast({
      title: isSignUp ? 'Регистрация успешна' : 'Вход выполнен',
      description: 'Добро пожаловать!',
    });
    
    navigate('/');
  };

  const toggleMode = () => {
    setIsSignUp(!isSignUp);
    form.reset();
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold text-center">
            {isSignUp ? 'Регистрация' : 'Вход'}
          </CardTitle>
          <CardDescription className="text-center">
            {isSignUp 
              ? 'Создайте аккаунт для начала работы' 
              : 'Войдите в свой аккаунт'}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Имя пользователя</FormLabel>
                    <FormControl>
                      <Input 
                        placeholder="Введите имя" 
                        autoComplete="username"
                        {...field} 
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Пароль</FormLabel>
                    <FormControl>
                      <Input 
                        type="password" 
                        placeholder="Введите пароль"
                        autoComplete={isSignUp ? 'new-password' : 'current-password'}
                        {...field} 
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button 
                type="submit" 
                className="w-full" 
                disabled={isSubmitting}
              >
                {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isSignUp ? 'Зарегистрироваться' : 'Войти'}
              </Button>
            </form>
          </Form>
          <div className="mt-4 text-center text-sm">
            <span className="text-muted-foreground">
              {isSignUp ? 'Уже есть аккаунт?' : 'Нет аккаунта?'}
            </span>{' '}
            <button
              type="button"
              onClick={toggleMode}
              className="text-primary hover:underline font-medium"
            >
              {isSignUp ? 'Войти' : 'Зарегистрироваться'}
            </button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
