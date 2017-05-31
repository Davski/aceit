-- MySQL dump 10.13  Distrib 5.7.18, for Win32 (AMD64)
--
-- Host: localhost    Database: aceit
-- ------------------------------------------------------
-- Server version	5.7.18-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `alt`
--

DROP TABLE IF EXISTS `alt`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `alt` (
  `altID` int(11) NOT NULL AUTO_INCREMENT,
  `alt1` char(200) NOT NULL,
  `alt2` char(200) NOT NULL,
  `alt3` char(200) NOT NULL,
  `courseID` int(4) NOT NULL,
  `moduleID` int(4) NOT NULL,
  `questionID` int(4) NOT NULL,
  PRIMARY KEY (`altID`),
  KEY `questionID` (`questionID`),
  KEY `moduleID` (`moduleID`),
  KEY `courseID` (`courseID`),
  CONSTRAINT `alt_ibfk_1` FOREIGN KEY (`questionID`) REFERENCES `question` (`questionID`),
  CONSTRAINT `alt_ibfk_2` FOREIGN KEY (`moduleID`) REFERENCES `module` (`moduleID`),
  CONSTRAINT `alt_ibfk_3` FOREIGN KEY (`courseID`) REFERENCES `course` (`courseID`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `alt`
--

LOCK TABLES `alt` WRITE;
/*!40000 ALTER TABLE `alt` DISABLE KEYS */;
INSERT INTO `alt` VALUES (1,'Process Control Block','Powered Circuit Board','Process Control Board',1,1,1),(4,'PCB is a layer between the user and the operative system, restricting the users access of kernel functions','PCB is a data structure in the operative system kernel conaining the information needed to manage a particular process','PCB is a storage type used in kernel programing',1,1,2),(5,'A process is a set of instructions which is in human readable format and a program is a process loaded into memory and executing or waiting.','A program is a set of instructions which is in human readable format and a process is program loaded into memory and executing or waiting.','A program is an active entity and needs resources such as CPU time, memory etc. to execute and a process is a passive entity stored on secondary storage.',1,1,3),(6,'New, running, waiting, ready, terminated, old.','In, running, waiting, ready, executed.','New, running, waiting, ready, terminated.',1,1,4),(7,'It creates an identical child process from a parent process, but they don?t run concurrently.','It kills an identical child process from it?s parent process.','It creates an identical child process from a parent process and they run concurrently',1,2,5),(8,'2','4','1',1,2,6),(9,'Returns a value of 1 to the child process and returns the process ID of the child process to the parent process.','Returns a value 0 to the parent process and returns the process ID of parent to the CPU.','Returns a value of 0 to the child process and returns the process ID of the child process to the parent process.',1,2,7),(10,'Waits for a child process to finish','Waits for a parent process to finish','Waits for both a parent and its child process to finish',1,2,8),(11,'The PID of the executing child process','The PID of the terminated child process.','The PID of the waiting child process',1,2,9),(12,'It is a special kernel object and it is a a pair of descriptors connected together. It is a simple LIFO communication channel.','It is a special kernel object and it is a pair of descriptors connected together. It is a simple FIFO communication channel.','It is a special OS object and it is a pair of descriptors connected together. It is a simple LIFO communication channel.',1,2,10),(13,'System call pipe()','System call createPipe()','System call forkPipe()',1,2,11),(14,'Program continues to next process','SIGPIPE signal, will cause process to terminate','Reading will be blocked()',1,2,12),(17,'SJF','FCFS','Round Robin',1,3,14),(18,'Process is preempted and added to the end of the waiting queue','Process is terminated','Process is preempted and added to the end of the ready queue',1,3,15),(19,'Knowing the length of the next CPU request','Hard to find a fitting time quantum','The convoy effect',1,3,16),(20,'Throughput time is numbers of processes that wait in the ready queue per time unit','Total time a process has spent in every state after it has been brought into the memory.','Numbers of processes that complete their execution per time unit',1,3,18),(21,'Amount of time to execute a particular process','Numbers of processes that complete their execution per time unit','Total time a process has spent in every state after it has been brought into the memory.',1,3,19),(22,'Minimizing the average waiting time','Preventing starvation','To make the the whole process as easy as possible',1,3,20),(23,'They share code, data, stack, and register, but not open files','They share code, stack and register, but not data and open files.','They share code, data and open files, but not stack and register',1,3,21),(24,'Pipes and shared memory','Pipes and phones','Telekinesis and shared memory',1,1,22);
/*!40000 ALTER TABLE `alt` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `answer`
--

DROP TABLE IF EXISTS `answer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `answer` (
  `answerID` int(11) NOT NULL AUTO_INCREMENT,
  `answer` char(100) NOT NULL,
  `courseID` int(4) NOT NULL,
  `moduleID` int(4) NOT NULL,
  `questionID` int(4) NOT NULL,
  PRIMARY KEY (`answerID`),
  KEY `questionID` (`questionID`),
  KEY `moduleID` (`moduleID`),
  KEY `courseID` (`courseID`),
  CONSTRAINT `answer_ibfk_1` FOREIGN KEY (`questionID`) REFERENCES `question` (`questionID`),
  CONSTRAINT `answer_ibfk_2` FOREIGN KEY (`moduleID`) REFERENCES `module` (`moduleID`),
  CONSTRAINT `answer_ibfk_3` FOREIGN KEY (`courseID`) REFERENCES `course` (`courseID`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `answer`
--

LOCK TABLES `answer` WRITE;
/*!40000 ALTER TABLE `answer` DISABLE KEYS */;
INSERT INTO `answer` VALUES (1,'alt1',1,1,1),(2,'alt2',1,1,2),(3,'alt2',1,1,3),(4,'alt3',1,1,4),(5,'alt3',1,2,5),(6,'alt1',1,2,6),(7,'alt3',1,2,7),(8,'alt1',1,2,8),(9,'alt2',1,2,9),(10,'alt2',1,2,10),(11,'alt1',1,2,11),(12,'alt3',1,2,12),(16,'alt3',1,3,18),(17,'alt1',1,3,19),(18,'alt1',1,3,20),(19,'alt3',1,3,21),(20,'alt1',1,1,22),(21,'alt2',1,3,14),(22,'alt3',1,3,15),(24,'alt1',1,3,16);
/*!40000 ALTER TABLE `answer` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course`
--

DROP TABLE IF EXISTS `course`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `course` (
  `courseID` int(11) NOT NULL AUTO_INCREMENT,
  `course` char(20) NOT NULL,
  PRIMARY KEY (`courseID`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course`
--

LOCK TABLES `course` WRITE;
/*!40000 ALTER TABLE `course` DISABLE KEYS */;
INSERT INTO `course` VALUES (1,'Ospp');
/*!40000 ALTER TABLE `course` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `module`
--

DROP TABLE IF EXISTS `module`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `module` (
  `moduleID` int(11) NOT NULL AUTO_INCREMENT,
  `module` char(20) NOT NULL,
  `courseID` int(4) NOT NULL,
  PRIMARY KEY (`moduleID`),
  KEY `courseID` (`courseID`),
  CONSTRAINT `module_ibfk_1` FOREIGN KEY (`courseID`) REFERENCES `course` (`courseID`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `module`
--

LOCK TABLES `module` WRITE;
/*!40000 ALTER TABLE `module` DISABLE KEYS */;
INSERT INTO `module` VALUES (1,'Module0',1),(2,'Module1',1),(3,'Module2',1);
/*!40000 ALTER TABLE `module` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `question`
--

DROP TABLE IF EXISTS `question`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `question` (
  `questionID` int(11) NOT NULL AUTO_INCREMENT,
  `question` char(200) NOT NULL,
  `courseID` int(4) NOT NULL,
  `moduleID` int(4) NOT NULL,
  PRIMARY KEY (`questionID`),
  KEY `moduleID` (`moduleID`),
  KEY `courseID` (`courseID`),
  CONSTRAINT `question_ibfk_1` FOREIGN KEY (`moduleID`) REFERENCES `module` (`moduleID`),
  CONSTRAINT `question_ibfk_2` FOREIGN KEY (`courseID`) REFERENCES `course` (`courseID`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `question`
--

LOCK TABLES `question` WRITE;
/*!40000 ALTER TABLE `question` DISABLE KEYS */;
INSERT INTO `question` VALUES (1,'What does the acronym PCB stand for?',1,1),(2,'What is the purpose of the PCB?',1,1),(3,'Which of the following statements is true?',1,1),(4,'What are the states that a process can be in?',1,1),(5,'What is true about the system call fork()?',1,2),(6,'How many times does fork() return?',1,2),(7,'What are the possible return values of fork()?',1,2),(8,'What is true about the wait() system call?',1,2),(9,'What are the possible return values of wait()?',1,2),(10,'What is true about a pipe?',1,2),(11,'How do we create a pipe?',1,2),(12,'What happens if we read from a pipe with no data?',1,2),(14,'Which alternative fit the following description: Simple, straight to point, FIFO',1,3),(15,'What happens after the given time for a process has elapsed when scheduled in Round Robin?',1,3),(16,'Which one is the disadvantage of the SJF?',1,3),(18,'What is a throughput time?',1,3),(19,'What is a turnaround time?',1,3),(20,'What is the overall purpose of SJF?',1,3),(21,'Which of the following statements about threads is true?',1,3),(22,'What are two different categories of inter process communication?',1,1);
/*!40000 ALTER TABLE `question` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `userID` int(11) NOT NULL AUTO_INCREMENT,
  `name` char(30) NOT NULL,
  `password` binary(60) DEFAULT NULL,
  PRIMARY KEY (`userID`)
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'Gustav','$2a$10$UcwlW8CDgAKo859lmUg8Euqr4zClKLhem0lcrTfmrCm7noFwIgFiu'),(37,'Ida','$2a$10$siUUDt9Wz1evXcvl3mccpeNCR6WCZKIZxgWnTAGdwEUD3PJEWrzDC'),(42,'q','$2a$10$am6wSUplXGHfT8ciQhtQL.cA/0CdUOSM/alZ9K68adgIokjcVtAV.'),(43,'a','$2a$10$xjmJd2ojPlZG4NouQAWRk.jBu./Uis9KcU1pJDLRGIaUwR.DGRhre');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-05-26 10:21:32
