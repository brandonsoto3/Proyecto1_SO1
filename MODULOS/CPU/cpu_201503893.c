#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>   
#include <linux/fs.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h> 
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>
#include <linux/sched.h>        
#include <linux/sched/signal.h> 

//DEFINIENDO EL TAMAÑO DEL BUFFER
#define BUFSIZE  	150

//DEFINIENDO INFORMACIÓN ACERCA DEL MODULO
MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Modulo que escribe informacion acerca del CPU");
MODULE_AUTHOR("BRANDON JAVIER SOTO CASTAÑEDA 201503893");

struct sysinfo inf;

int contador = 0;
long total_memoria = 0;


void imprimir_estado(struct task_struct *task_next, struct seq_file * m)
{

    
        seq_printf(m,task_next);

    if(task_next->state == -1){
        seq_printf(m,"      \"Estado\":\"NO EJECUTABLE\"");
    }else if(task_next->state == 0){
        seq_printf(m,"      \"Estado\":\"EJECUCION\"");
    }else if(task_next->state == 1){
        seq_printf(m,"      \"Estado\":\"SUSPENDIDO\""); //INTERRUPTIBLE
    }else if(task_next->state == 2){
        seq_printf(m,"      \"Estado\":\"SUSPENDIDO\""); //UNINTERRUPTIBLE
    }else if(task_next->state == 4){ 
        seq_printf(m,"      \"Estado\":\"ZOMBIE\""); //STOPEED
    }else if(task_next->state == 8){ 
        seq_printf(m,"      \"Estado\":\"DETENIDO\"");//TRACED
    }else{
        seq_printf(m,"      \"Estado\":\"DESCONOCIDO\"");
    }
}


void recorrer_tareas(struct task_struct *task, struct seq_file * m, int num)
{
  struct task_struct *task_next;
  struct list_head *list;
  
  int tmp = contador;
  contador = 1;
  
  int cant = 0;
  list_for_each(list, &task->children) {
    task_next = list_entry(list, struct task_struct, sibling);
    if(task_next == 0){
        continue;
    }
    cant = cant+1;
  }
  
  list_for_each(list, &task->children) {
    task_next = list_entry(list, struct task_struct, sibling);

    if(task_next == 0){
        continue;
    }
    seq_printf(m,"      {\n");    
    seq_printf(m,"      \"PID\":%d,\n",task_next->pid);	
    seq_printf(m,"      \"Nombre\":\"%s\",\n",task_next->comm);
    seq_printf(m,"      \"Usuario\":\"%d\",\n",task_next->cred->uid);
    
    //IMPRIMIMOS EL ESTADO
    imprimir_estado(task_next,m);

    seq_printf(m,",\n");

    //Porcentaje de la memoria RAM 
    long memoria_utilizada = 0;
    if(task_next->mm){
        memoria_utilizada = (task_next->mm->total_vm <<(PAGE_SHIFT -10));
        memoria_utilizada = (memoria_utilizada/58603);
    }

    //memoria_utilizada*100/total_memoria
    seq_printf(m,"      \"Porcentaje_RAM\": %i",memoria_utilizada);
    seq_printf(m,",\n");

    //PPID
    seq_printf(m,"      \"PPID\":%d\n",num);
    
    if(cant > 1){
        seq_printf(m,"      },\n");
        cant = cant - 1;
    }else {
        if(num == 2){
            seq_printf(m,"      }\n");    
        }else{
            seq_printf(m,"      },\n");
        }
    }
    //EN ESTA PARTE VOLVEMOS A LLAMAR A LA FUNCION PARA RECORRER A LOS HIJOS
    recorrer_tareas(task_next, m, task_next->pid); 
  }
  contador = tmp;
}


static int write_file(struct seq_file * m, void *v) {	
    contador = 1;
    si_meminfo(&inf);
    total_memoria 	= (inf.totalram * 4)/1024;    
    seq_printf(m,"{\n");
    seq_printf(m,"      \"lista_procesos\":[\n");
    recorrer_tareas(&init_task, m, 0);
    seq_printf(m,"      ]\n");
    seq_printf(m,"}\n");
    return 0;
}

static int open(struct inode *inode, struct  file *file) {
  return single_open(file, write_file, NULL);
}

static struct file_operations operaciones =
{    
    .open = open,
    .read = seq_read
};


static int iniciar(void)
{
    proc_create("cpu_201503893", 0, NULL, &operaciones);
    printk(KERN_INFO "BRANDON JAVIER SOTO CASTANEDA\n");
    return 0;
}
 
static void salir(void)
{
    remove_proc_entry("cpu_201503893", NULL);
    printk(KERN_INFO "Diciembre 2020\n");
}
 
module_init(iniciar);
module_exit(salir); 